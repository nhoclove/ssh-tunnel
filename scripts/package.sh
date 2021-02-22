#!/usr/bin/env bash

tool_name=$1
if [[ -z "$tool_name" ]]; then
    echo "usage: $0 <package-name>"
    exit 1
fi

WORKDIR="." #"$(cd "$(dirname "$0")" && pwd)"
PLATFORMS=("windows/amd64" "linux/amd64")

for platform in "${PLATFORMS[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    binary_name=$tool_name'_'$GOOS'_'$GOARCH
    package_name=$tool_name'_'$GOOS'_'$GOARCH
    if [ $GOOS = "windows" ]; then
        binary_name+='.exe'
    fi

    # create a folder for each platform
    tmp="$WORKDIR/dist/$package_name"
    mkdir -p -- "$tmp"
    cp -r "$WORKDIR/config/" "$tmp"

    env GOOS=$GOOS GOARCH=$GOARCH go build -mod=vendor -o "$tmp/$binary_name" "$WORKDIR"
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    # zip package
    zip -r "$tmp.zip" "$tmp"
    rm -rf "$tmp"
done
