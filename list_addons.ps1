<#
.SYNOPSIS
    Lists all installed Garry's Mod addons with their names and information from Steam Workshop.
.DESCRIPTION
    This script scans the addons/0/out directory and fetches information for each addon
    using the Steam Workshop API.
#>

$gmod_dir = "C:\Local\Garrys Mod\"
$out_dir = "$gmod_dir\garrysmod\addons\0\out"

# Steam Workshop API key (replace with your own or leave empty for limited access)
$steam_api_key = ""

# Check if out directory exists
if (-not (Test-Path $out_dir)) {
    Write-Error "Out directory not found at $out_dir"
    exit 1
}

# Get all addon directories in out folder
$addon_folders = Get-ChildItem -Path $out_dir -Directory

if ($addon_folders.Count -eq 0) {
    Write-Host "No addons found in $out_dir"
    exit
}

Write-Host "Installed Addons (from Steam Workshop):"
Write-Host "======================================"
Write-Host ""

foreach ($folder in $addon_folders) {
    $addon_id = $folder.Name

    # Try to fetch addon info from Steam Workshop API
    try {
        $api_url = "https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/"

        # Create the request body as form-urlencoded
        $requestBody = "itemcount=1"
        $requestBody += "&publishedfileids[0]=$addon_id"

        if (-not [string]::IsNullOrEmpty($steam_api_key)) {
            $requestBody += "&key=$steam_api_key"
        }

        $response = Invoke-RestMethod -Uri $api_url -Method Post -Body $requestBody -ContentType "application/x-www-form-urlencoded"

        if ($response.response.publishedfiledetails -and $response.response.publishedfiledetails.Count -gt 0) {
            $details = $response.response.publishedfiledetails[0]

            Write-Host "ID: $addon_id"
            Write-Host "Title: $($details.title)"
            Write-Host "Creator: $($details.creator)"
            Write-Host "Time Created: $([datetimeoffset]::FromUnixTimeSeconds($details.time_created).DateTime.ToString('yyyy-MM-dd HH:mm:ss'))"
            Write-Host "Time Updated: $([datetimeoffset]::FromUnixTimeSeconds($details.time_updated).DateTime.ToString('yyyy-MM-dd HH:mm:ss'))"
            Write-Host "Views: $($details.views)"
            Write-Host "Subscriptions: $($details.subscriptions)"
            Write-Host "Favorites: $($details.favorited)"
            Write-Host "Tags: $($details.tags -join ', ')"
            Write-Host "Description: $($details.description)"
            Write-Host ""
        } else {
            Write-Host "ID: $addon_id"
            Write-Host "No information available from Steam Workshop"
            Write-Host ""
        }
    }
    catch {
        Write-Host "ID: $addon_id"
        Write-Host "Error fetching information: $_"
        Write-Host ""
    }
}
