#!/bin/bash

# Create build directory if it doesn't exist
mkdir -p build

# Build based on platform argument
case "$1" in
    "windows")
        GOOS=windows GOARCH=amd64 go build -o build/gmod-addon-manager-windows.exe
        ;;
    "linux")
        GOOS=linux GOARCH=amd64 go build -o build/gmod-addon-manager-linux
        ;;
    *)
        echo "Usage: $0 [linux|windows]"
        exit 1
        ;;
esac
