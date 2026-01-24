# GMod Addon Manager

A terminal-based application for downloading, installing, and managing Garry's Mod addons from the Steam Workshop.

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

