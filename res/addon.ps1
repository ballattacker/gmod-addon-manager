# usage: addon.ps1 -id 123

param (
    [string]$id
)

$gmod_dir = "C:\Local\Garrys Mod\"
$download_dir = "C:\Users\ADMIN\AppData\Local\Microsoft\WinGet\Packages\Valve.SteamCMD_Microsoft.Winget.Source_8wekyb3d8bbwe\steamapps\workshop\content\4000\"
$addon_dir = "$gmod_dir\garrysmod\addons"
$tmp_dir = "$addon_dir\0\tmp"
$out_dir = "$addon_dir\0\out"

New-Item -ItemType Directory -Force -Path $tmp_dir
New-Item -ItemType Directory -Force -Path $out_dir

& "steamcmd.exe" +login anonymous +workshop_download_item 4000 $id +quit

# Find the downloaded file (either .gma or _legacy.bin)
$downloaded_file = Get-ChildItem -Path "$download_dir\$id" -File | Select-Object -First 1
$file_path = "$download_dir\$id\$($downloaded_file.Name)"
$file_name = $downloaded_file.Name

New-Item -ItemType Directory -Force -Path "$tmp_dir\$id"

# Handle .bin file (extract and rename to .gma)
if ($file_name -like "*_legacy.bin") {
    # Extract the .bin file (assuming it's a zip)
    Expand-Archive -Path $file_path -DestinationPath "$tmp_dir\$id" -Force
    # Find the extracted .gma file
    $extracted = Get-ChildItem -Path "$tmp_dir\$id" | Select-Object -First 1
    # Move the extracted .gma to tmp directory
    Move-Item "$tmp_dir\$id\$($extracted.Name)" "$tmp_dir\$id\$id.gma"
} elseif ($file_name -like "*.gma") {
    # Move the extracted .gma to tmp directory
    Move-Item $file_path "$tmp_dir\$id\$id.gma"
} else {
    # error
}

# Remove the original downloaded folder
Remove-Item "$download_dir\$id" -Recurse -Force

# Execute GMAD tool
& "$gmod_dir\bin\gmad.exe" "$tmp_dir\$id\$id.gma"

# Move the extracted addon to the addons directory with new name
Move-Item -Path "$tmp_dir\$id\$id" -Destination "$out_dir\$id"

# Clean up tmp directory
Remove-Item "$tmp_dir\$id" -Recurse -Force

# Symlink to enable addon
New-Item -ItemType SymbolicLink -Path "$addon_dir\$id" -Value "$out_dir\$id"
