# 🎮 GMod Addon Manager

[![Release](https://img.shields.io/github/v/release/yourusername/gmod-addon-manager?style=flat-square)](https://github.com/yourusername/gmod-addon-manager/releases)
[![License](https://img.shields.io/github/license/yourusername/gmod-addon-manager?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/yourusername/gmod-addon-manager?style=flat-square)](go.mod)

A terminal-based application for downloading, installing, and managing Garry's Mod addons from the Steam Workshop.

## Features

- 📥 Download and install addons from Steam Workshop
- ⚡ Enable/disable installed addons
- 🗑️ Remove addons (including files)
- 📋 View information about installed addons
- 💾 Cache addon information to reduce API calls
- 🖥️ Both TUI (Terminal User Interface) and CLI modes

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/yourusername/gmod-addon-manager/releases).

### Building from Source

#### Prerequisites

- Go 1.24+
- SteamCMD installed and in your PATH
- Garry's Mod installed

#### Steps

1. Clone this repository:
   ```shell
   git clone https://github.com/yourusername/gmod-addon-manager.git
   cd gmod-addon-manager
   ```

2. Build the application:
   ```shell
   go build
   ```

3. The binary will be created in the current directory

## Usage

### TUI Mode

Run the application without arguments to launch the interactive TUI:

```shell
gmod-addon-manager
```

### CLI Mode

The application also supports command-line usage:

```shell
gmod-addon-manager [command]
```

Available commands:

- `get [addon-id]` - Download and install an addon
- `enable [addon-id]` - Enable an installed addon
- `disable [addon-id]` - Disable an installed addon
- `remove [addon-id]` - Remove an addon
- `list` - List all installed addons
- `info [addon-id]` - Show information about an addon
- `config` - Show current configuration

## Configuration

On first run, the application will create a default configuration file at:

- Windows: `%APPDATA%\gmod-addon-manager\gmod-addon-manager.json`
- Linux/macOS: `~/.config/gmod-addon-manager/gmod-addon-manager.json`

You can edit this file to customize paths and settings.

## Releases

Check out the [Releases page](https://github.com/yourusername/gmod-addon-manager/releases) for pre-built binaries and changelog information.

## License

[MIT License](LICENSE)
