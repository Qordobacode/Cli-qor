#!/usr/bin/env bash

platforms=("windows/amd64" "darwin/amd64" "linux/amd64" "linux/arm64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='qor-'$GOOS'-'$GOARCH
    if [[ $GOOS = "windows" ]]; then
        output_name+='.exe'
    fi
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X 'github.com/qordobacode/cli-v2/cmd/info.VersionFlag=v0.8.4'" -o $output_name
    if [[ $? -ne 0 ]]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done