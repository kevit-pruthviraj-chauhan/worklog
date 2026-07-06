# GitHub Actions Auto-Release Setup

## Overview
The worklog project now has automatic binary building and releasing via GitHub Actions. When you push a version tag, GitHub Actions will automatically:
1. Build binaries for all supported platforms (Linux/Darwin, amd64/arm64)
2. Create a GitHub release
3. Upload all binaries to the release

## One-Time Setup

### Step 1: Push Code to GitHub
Make sure your code is pushed to GitHub:
```bash
cd /home/kevit/checkin_checkout/worklog
git remote add origin https://github.com/kevit-pruthviraj-chauhan/worklog.git  # if not already added
git push -u origin main  # or your branch
```

### Step 2: Enable GitHub Actions
The workflow file is already in `.github/workflows/release.yml`. GitHub Actions should automatically detect it.

## How to Create a Release

### Step 1: Create a Git Tag
```bash
# Create a tag (replace v1.0.0 with your version)
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

### Step 2: Monitor the Build
Go to your GitHub repository and click on the **Actions** tab to see the workflow running.

### Step 3: Release is Created Automatically
Once the workflow completes, a new release will be available at:
`https://github.com/kevit-pruthviraj-chauhan/worklog/releases`

The release will contain:
- `worklog-linux-amd64` - Linux 64-bit
- `worklog-linux-arm64` - Linux ARM64 (Raspberry Pi, etc.)
- `worklog-darwin-amd64` - macOS Intel
- `worklog-darwin-arm64` - macOS Apple Silicon

## Using the Auto-Update Feature

Once you have created a release with binaries, users can update with:
```bash
worklog update
```

This will automatically download and install the latest version.

## Example Workflow

1. Make changes to your code
2. Test locally: `make build` and test the binary
3. Commit your changes: `git commit -am "your message"`
4. Create a tag: `git tag v1.0.1`
5. Push: `git push origin main && git push origin v1.0.1`
6. GitHub Actions automatically builds and releases!
7. Users can now run `worklog update` to get the latest version

## Troubleshooting

### Release not appearing
- Check the **Actions** tab on GitHub - the workflow may still be running
- Ensure your tag follows the pattern `v*` (e.g., v1.0.0, v1.1.0)
- Check workflow logs for any build errors

### Update command still fails with 404
- Make sure at least one release has been published
- Verify the GitHub repository URL is correct
- Wait a few minutes after pushing the tag for the release to be created

## Manual Release (Without GitHub Actions)

If you need to manually create a release:

```bash
# Build all binaries
make build-all

# Go to GitHub releases page and create a new release
# Upload all the worklog-* files
```

## Variables in the Workflow

The workflow uses the following:
- **Go version**: 1.25 (specified in setup-go)
- **Platforms**: Linux (amd64, arm64) and Darwin/macOS (amd64, arm64)
- **Token**: Uses GitHub's built-in GITHUB_TOKEN (no setup needed)

To add more platforms or change Go version, edit `.github/workflows/release.yml`
