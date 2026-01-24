#!/bin/bash

id=$1
name=${2:?}
steamcmd +login anonymous +workshop_download_item 4000 "$id" +quit
gmod_dir="$HOME/.user/loc/gmod/GarrysMod"
download_dir="$HOME/.steam/SteamApps/workshop/content/4000"
addon_dir="$gmod_dir/game/garrysmod/addons"
tmp_dir="$HOME/.user/tmp/gmod-addon"
name=$(perl -Mutf8 -pe 's/[^\w]+/_/g; s/_+/_/g; s/^_|_$//g; $_ = lc' <<< "$name")
mkdir -p "$tmp_dir"
mv "$download_dir/$id"/*.gma "$tmp_dir/$id.gma"
rm -r "${download_dir:?}/${id:?}"
"$gmod_dir"/game/bin/linux64/gmad "$tmp_dir/$id.gma"
rm "$tmp_dir/$id.gma"
mv "$tmp_dir/$id" "$addon_dir/$id-$name"
