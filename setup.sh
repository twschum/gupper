#!/bin/bash
# builds the app, server, and some packages
# builds the docker container for the server, if possible

go install ./...
bash build_packages.sh 1.0.0

if [ $(command -v docker) ]; then
    docker build -t gupper-server -f Dockerfile .
else
    echo "docker command not available, skipping build"
fi
