<#
.SYNOPSIS
    Lists all installed Garry's Mod addons with their names and information.
.DESCRIPTION
    This script scans the Garry's Mod addons directory and displays information
    about each installed addon by reading their addon.json files.
#>

$gmod_dir = "C:\Local\Garrys Mod\"
$addon_dir = "$gmod_dir\garrysmod\addons"

# Check if addons directory exists
if (-not (Test-Path $addon_dir)) {
    Write-Error "Addons directory not found at $addon_dir"
    exit 1
}

# Get all addon directories (excluding symlinks and special directories)
$addon_folders = Get-ChildItem -Path $addon_dir -Directory |
                 Where-Object { $_.Name -notmatch '^0$' } |
                 Where-Object { -not $_.Attributes.HasFlag([System.IO.FileAttributes]::ReparsePoint) }

if ($addon_folders.Count -eq 0) {
    Write-Host "No addons found in $addon_dir"
    exit
}

Write-Host "Installed Addons:"
Write-Host "================"
Write-Host ""

foreach ($folder in $addon_folders) {
    $addon_path = Join-Path -Path $addon_dir -ChildPath $folder.Name
    $addon_json = Join-Path -Path $addon_path -ChildPath "addon.json"

    # Try to read addon.json
    if (Test-Path $addon_json) {
        try {
            $addon_data = Get-Content -Path $addon_json -Raw | ConvertFrom-Json

            Write-Host "ID: $($folder.Name)"
            Write-Host "Title: $($addon_data.title)"
            Write-Host "Author: $($addon_data.author)"
            Write-Host "Version: $($addon_data.version)"
            Write-Host "Description: $($addon_data.description)"
            Write-Host "Tags: $($addon_data.tags -join ', ')"
            Write-Host ""
        }
        catch {
            Write-Host "ID: $($folder.Name)"
            Write-Host "Error reading addon.json: $_"
            Write-Host ""
        }
    }
    else {
        Write-Host "ID: $($folder.Name)"
        Write-Host "No addon.json found"
        Write-Host ""
    }
}
