# GMod Addon Manager

A terminal-based application for downloading, installing, and managing Garry's
Mod addons from the Steam Workshop.

## Features

- Download and install addons from Steam Workshop
- Enable/disable installed addons
- Remove addons (including files)
- View information about installed addons
- Cache addon information to reduce API calls
- Both TUI (Terminal User Interface) and CLI modes

## Installation

### Prerequisites

- Go 1.20+
- SteamCMD installed and in your PATH
- Garry's Mod installed

### Building from source

1. Clone this repository
2. Run `go build` to build the application
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

## License

[MIT License](LICENSE)
