#!/bin/bash
################################################################################

version=$1
if [[ -z "$version" ]]; then
    echo "usage: $0 <version> [package_dir]"
    exit 1
fi
package_dir=$2
if [[ -z "$package_dir" ]]; then
    package_dir=packages/
fi
package_name=app
package=cmd/$package_name/$package_name.go

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$version'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.BuildVersion=$version" -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution…'
        exit 1
    fi
    echo $output_name
    mv $output_name $package_dir
done
