#!/bin/bash
echo "=== Building Myst for production..."

echo "=== Cleaning up old files..."
rm -r "build" > /dev/null 2>&1

if [[ -z "${GOOS}" ]]; then
    export GOOS=linux
fi
if [[ -z "${GOARCH}" ]]; then
    export GOARCH=amd64
fi

echo "--- Platform: $GOOS"
echo "--- Architecture: $GOARCH"

echo "=== Building server..."

if [ $GOOS = windows ]; then
    go build -o build/server/server-${GOOS}-${GOARCH}.exe cmd/server/main.go
else
    go build -o build/server/server-${GOOS}-${GOARCH} cmd/server/main.go
fi

echo "=== Building UI..."

rm -rf ./static

cd ui
npm run build
cd ..

rm -rf ./cmd/client/static
cp -r static ./cmd/client/static

if [ $GOOS = windows ]; then
    go build -o build/client/client-${GOOS}-${GOARCH}.exe cmd/client/main.go
else
    go build -o build/client/client-${GOOS}-${GOARCH} cmd/client/main.go
fi

#echo "Compiling frontend bundle..."
#cd vue > /dev/null 2>&1
#npm install > /dev/null 2>&1
#npm run build > /dev/null 2>&1
#cd .. > /dev/null 2>&1
#echo "Frontend bundle compiled."
#echo "Copying assets..."
#cp -r "static" "build/static"  > /dev/null 2>&1
#cp -r ".env."* "build" > /dev/null 2>&1
#cp -r "assets" "build/assets" > /dev/null 2>&1
#echo "Build Successful."
#sleep 3
