#!/bin/bash
echo "Building for production..."
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
    go build -o build/server-${GOOS}-${GOARCH}.exe go/main.go
else
    go build -o build/server-${GOOS}-${GOARCH} go/main.go
fi
echo "Server binary compiled."
echo "Compiling frontend bundle..."
cd vue > /dev/null 2>&1
npm install > /dev/null 2>&1
npm run build > /dev/null 2>&1
cd .. > /dev/null 2>&1
echo "Frontend bundle compiled."
echo "Copying assets..."
cp -r "static" "build/static"  > /dev/null 2>&1
cp -r ".env."* "build" > /dev/null 2>&1
cp -r "assets" "build/assets" > /dev/null 2>&1
echo "Build Successful."
sleep 3
