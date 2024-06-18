#!/bin/bash

mkdir -p bin

# Build the application
go build -o bin/tasktango .

echo "Build completed"
