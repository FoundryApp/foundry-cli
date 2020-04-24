#!/bin/sh

env GOOS=darwin GOARCH=amd64 go build -o ./build/foundry-macos-x86_64
env GOOS=linux GOARCH=amd64 go build -o ./build/foundry-linux-x86_64
env GOOS=linux GOARCH=arm64 go build -o ./build/foundry-linux-arm64