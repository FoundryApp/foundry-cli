#!/bin/sh

env GOOS=linux GOARCH=amd64 go build -tags debug -o "./build/foundry-debug-linux" .