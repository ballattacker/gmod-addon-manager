# usage: addon.ps1 -id 123

param (
    [string]$id
)

& "steamcmd.exe" +login anonymous +workshop_download_item 4000 $id +quit

$gmod_dir = "C:\Local\Garrys Mod\"
$download_dir = "C:\Users\ADMIN\AppData\Local\Microsoft\WinGet\Packages\Valve.SteamCMD_Microsoft.Winget.Source_8wekyb3d8bbwe\steamapps\workshop\content\4000\"
$addon_dir = "$gmod_dir\garrysmod\addons"
$tmp_dir = "$addon_dir\0\tmp"
$out_dir = "$addon_dir\0\out"

New-Item -ItemType Directory -Force -Path $tmp_dir
New-Item -ItemType Directory -Force -Path $out_dir

# Find the downloaded file (either .gma or _legacy.bin)
$downloaded_file = Get-ChildItem -Path "$download_dir\$id" -File | Select-Object -First 1
$file_path = "$download_dir\$id\$($downloaded_file.Name)"
$file_name = $downloaded_file.Name

# Handle .bin file (extract and rename to .gma)
if ($file_name -like "*_legacy.bin") {
    # Extract the .bin file (assuming it's a zip)
    Expand-Archive -Path $file_path -DestinationPath "$tmp_dir\bin_extract" -Force
    # Find the extracted .gma file
    $extracted_gma = Get-ChildItem -Path "$tmp_dir\bin_extract" -Filter "*.gma" | Select-Object -First 1
    $gma_name = $extracted_gma.Name
    # Move the extracted .gma to tmp directory
    Move-Item "$tmp_dir\bin_extract\$gma_name" "$tmp_dir\$gma_name"
    # Clean up extraction directory
    Remove-Item "$tmp_dir\bin_extract" -Recurse -Force
} else {
    # It's already a .gma file
    $gma_name = $file_name
    # Move the .gma file to tmp directory
    Move-Item $file_path "$tmp_dir\$gma_name"
}

# Remove the original downloaded folder
Remove-Item "$download_dir\$id" -Recurse -Force

# Execute GMAD tool
& "$gmod_dir\bin\gmad.exe" "$tmp_dir\$gma_name"

# Remove the temporary GMA file
Remove-Item "$tmp_dir\$gma_name"

# Move the extracted addon to the addons directory with new name
Move-Item -Path "$tmp_dir\$id" -Destination "$out_dir\$id"

# Symlink to enable addon
New-Item -ItemType SymbolicLink -Path "$addon_dir\$id" -Value "$out_dir\$id"
