#!/usr/bin/env bash

# Build gpdb CLI
cd gpdb
rm -rf gpdb
env GOOS=linux GOARCH=amd64 go build -o gpdb

# Return back
cd ..

# Build datalab
cd datalab
rm -rf datalab
go build -o datalab

# Return back
cd ..


