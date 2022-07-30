#!/bin/bash
echo "=== Bootstrapping Myst..."

echo "=== Installing oapi-codegen..."
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

echo "=== Downloading node modules..."
cd ui && npm ci && cd ..

go generate ./...

echo "=== Downloading go modules..."
go mod download

echo "=== Done."
