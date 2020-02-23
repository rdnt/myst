#!/bin/bash
echo Building for production...
sleep 1
echo "Cleaning up old files..."
rm -r "build" > /dev/null 2>&1
echo "Formatting server packages..."
go fmt ./go/...
echo "Compiling server binary..."
if [[ -z "${GOOS}" ]]; then
    export GOOS=linux
fi
if [[ -z "${GOARCH}" ]]; then
    export GOARCH=amd64
fi
echo "  detected platform: $GOOS, architecture: $GOARCH"
if [ $GOOS = windows ]; then
    go build -o build/server-${GOOS}-${GOARCH}.exe go/server.go
else
    go build -o build/server-${GOOS}-${GOARCH} go/server.go
fi
echo "Server binary compiled."
echo "Compiling frontend bundle..."
npm --prefix vue/ run build > /dev/null 2>&1
echo "Frontend bundle compiled."
echo "Copying assets..."
cp -r "static" "build/static"
cp -r ".env."* "build"
cp -r "assets" "build/assets"
echo "Build Successful."
sleep 3
