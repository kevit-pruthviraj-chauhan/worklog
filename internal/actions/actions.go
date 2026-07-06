package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"worklog/internal/state"
)

func AddLog(s *state.State, msg string, t time.Time) {
	entry := fmt.Sprintf("%s → %s", t.Format("15:04:05"), msg)
	s.Logs = append(s.Logs, entry)
}

func nowLocal() time.Time {
	now := time.Now()
	return now.Local()
}

func ParseOptionalTime(arg string, ref time.Time) (time.Time, error) {
	if arg == "" {
		return nowLocal(), nil
	}

	parsed, err := time.Parse("15:04", arg)
	if err != nil {
		return time.Time{}, err
	}
	now := nowLocal()
	result := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, now.Location())
	if !ref.IsZero() && result.Before(ref) {
		candidate := result.Add(12 * time.Hour)
		if !candidate.Before(ref) && candidate.Sub(ref) < 24*time.Hour {
			return candidate, nil
		}
	}

	return result, nil
}

func Chkin(s *state.State, ts time.Time) error {
	switch s.Step {
	case 0:
		AddLog(s, "CHECKIN: Entry", ts)
		s.EntryTime = ts.Format(time.RFC3339)
		s.Step = 1
		return nil
	case 2:
		AddLog(s, "CHECKIN: Lunch End", ts)
		s.LunchEnd = ts.Format(time.RFC3339)
		s.Step = 3
		return nil
	default:
		return errors.New("cannot check in right now")
	}
}

func Chkout(s *state.State, ts time.Time) error {
	switch s.Step {
	case 1:
		AddLog(s, "CHKOUT: Lunch Start", ts)
		s.LunchStart = ts.Format(time.RFC3339)
		s.Step = 2
		return nil
	case 3:
		AddLog(s, "CHKOUT: Exit", ts)
		s.ExitTime = ts.Format(time.RFC3339)
		s.Step = 4
		return nil
	default:
		return errors.New("cannot check out right now")
	}
}

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

// Update downloads and installs the latest version of worklog
func Update() error {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Fetch latest release from GitHub
	resp, err := http.Get("https://api.github.com/repos/kevit-pruthviraj-chauhan/worklog/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release info: %w", err)
	}

	if release.TagName == "" {
		return errors.New("no releases found")
	}

	// Find the appropriate binary for this OS/arch
	osName := runtime.GOOS
	arch := runtime.GOARCH
	binaryName := fmt.Sprintf("worklog-%s-%s", osName, arch)

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == binaryName {
			downloadURL = asset.URL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s/%s in release %s", osName, arch, release.TagName)
	}

	// Download the binary
	fmt.Printf("Downloading worklog %s...\n", release.TagName)
	resp, err = http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "worklog-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write binary: %w", err)
	}
	tmpFile.Close()

	// Make it executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Try to replace the old binary with the new one
	// If permissions denied, provide helpful instructions
	oldExePath := exePath + ".old"
	if err := os.Rename(exePath, oldExePath); err != nil {
		// Check if it's a permission denied error
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied updating %s\n\nTo update, run:\n  sudo worklog update\n\nOr manually:\n  sudo mv %s %s", exePath, tmpFile.Name(), exePath)
		}
		return fmt.Errorf("failed to backup old binary: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), exePath); err != nil {
		// Restore the old binary if the new one fails to move
		os.Rename(oldExePath, exePath)
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied updating %s\n\nTo update, run:\n  sudo worklog update\n\nOr manually:\n  sudo mv %s %s", exePath, tmpFile.Name(), exePath)
		}
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Clean up the old binary
	os.Remove(oldExePath)

	fmt.Printf("Successfully updated to %s\n", release.TagName)
	fmt.Printf("Binary location: %s\n", exePath)
	return nil
}
