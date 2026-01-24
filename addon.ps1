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

# Move the downloaded GMA file to the temp directory
Move-Item "$download_dir\$id\*.gma" "$tmp_dir\$id.gma"
# Remove the original downloaded folder
Remove-Item "$download_dir\$id" -Recurse -Force

# Execute GMAD tool
& "$gmod_dir\bin\gmad.exe" "$tmp_dir\$id.gma"

# Remove the temporary GMA file
Remove-Item "$tmp_dir\$id.gma"

# Move the extracted addon to the addons directory with new name
Move-Item -Path "$tmp_dir\$id" -Destination "$out_dir\$id"

# Symlink to enable addon
New-Item -ItemType SymbolicLink -Path "$addon_dir\$id" -Value "$out_dir\$id"
