# Worklog Update Guide

## Overview
The worklog application now includes an automatic update feature that downloads and installs the latest version from GitHub releases.

## How to Use the Update Command

Simply run:
```bash
worklog update
```

The update command will:
1. Check GitHub for the latest release
2. Download the appropriate binary for your OS/architecture
3. Replace your current worklog binary
4. Display the location of the updated binary

## Setting Up GitHub Releases for Auto-Update

To enable the auto-update feature, you need to create GitHub releases with pre-built binaries. Follow these steps:

### Step 1: Build Binaries for All Platforms

```bash
make build-all
```

This creates binaries for:
- `worklog-linux-amd64`
- `worklog-linux-arm64`
- `worklog-darwin-amd64`
- `worklog-darwin-arm64`

### Step 2: Create a Git Tag and Push

```bash
git tag v1.0.0  # Use your version number
git push origin v1.0.0
```

### Step 3: Create a GitHub Release

1. Go to https://github.com/kevit-pruthviraj-chauhan/worklog/releases
2. Click "Create a new release"
3. Select your tag (e.g., v1.0.0)
4. Upload all the binary files created in Step 1
5. Publish the release

### Step 4: Users Can Now Update

Users can now run `worklog update` and it will automatically download and install the latest version.

## Installation Steps

### First Time Installation

1. Clone the repository:
```bash
git clone https://github.com/kevit-pruthviraj-chauhan/worklog.git
cd worklog
```

2. Build the binary:
```bash
make build
```

3. Install to /usr/local/bin:
```bash
make install
```

Or manually:
```bash
sudo mv worklog /usr/local/bin/
```

### Updating After Installation

```bash
worklog update
```

## Troubleshooting

### Update fails with "no binary found"
Make sure you have published a release on GitHub with the correct binary naming format:
- `worklog-linux-amd64`
- `worklog-linux-arm64`
- `worklog-darwin-amd64`
- `worklog-darwin-arm64`

### Permission denied errors
If you get permission errors when updating, ensure the worklog binary has proper permissions:
```bash
chmod +x /usr/local/bin/worklog
```

### Manual Update
If automatic update fails, you can always rebuild and reinstall manually:
```bash
cd ~/checkin_checkout/worklog
git pull
make install
```

## Automated Release Workflow (GitHub Actions)

You can automate binary building and release creation using GitHub Actions. Create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: [linux, darwin]
        arch: [amd64, arm64]
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      - run: GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} go build -o worklog-${{ matrix.os }}-${{ matrix.arch }} ./cmd
      - uses: actions/upload-artifact@v3
        with:
          name: worklog-${{ matrix.os }}-${{ matrix.arch }}
          path: worklog-${{ matrix.os }}-${{ matrix.arch }}

  create-release:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/download-artifact@v3
      - uses: softprops/action-gh-release@v1
        with:
          files: worklog-*
```

This will automatically build and upload binaries whenever you push a tag.
