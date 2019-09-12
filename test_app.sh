#!/bin/bash

# Run the go unit tests
echo "Starting unit tests"
go test ./...

unit=$1
if [ $? -ne 0 ] || [ "$unit" == "-u" ]; then
    exit
fi

# Run an end-to-end integration test on the app/server using the local docker host
echo "Setting up integration test"
port=8999
function cleanup {
    echo "cleaning up..."
    docker stop gupper-server-test
    docker rm gupper-server-test
    rm -rf test_pkgs
    exit 1
}

function check {
    if [ $? -ne 0 ]; then
        echo "Failure: $1"
        cleanup
    fi
}

# Make some packages
mkdir test_pkgs
bash build_version.sh 1.9.2 test_pkgs/
check "build packages"
bash build_version.sh 1.11 test_pkgs/
check "build packages"

# ephemeral test container that won't interfere with another one running default
docker run -d --name gupper-server-test -p $port:8080 -v $(pwd)/test_pkgs:/var/packages gupper-server
check "docker run"

# inital app is 1.0, will pick 1.11 > 1.9.2 and update to that version
go build -o app -ldflags "-X main.BuildVersion=1.0" cmd/app/app.go
check "app build"

# should update itself, check the output of --version
echo "Starting integration test"
./app --port :$port
check "app run"
output=$(./app -version)
if [ "$output" == "1.11.0" ]; then
    echo "Passed!"
else
    echo "app reports incorrect version: \"$output\", expected \"1.11.0\""
fi
cleanup
