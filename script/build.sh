#!/bin/bash
echo "=== Building Myst for production..."



echo "=== Building client..."




echo "=== Cleaning up old files..."
rm -r "build" > /dev/null 2>&1


echo "--- Detecting platform..."

if [[ -z "${GOOS}" ]]; then
    export GOOS=lpoinux
fi
if [[ -z "${GOARCH}" ]]; then
    export GOARCH=amd64
fi

echo "    --- platform:  $GOOS, architecture: $GOARCH"

echo "=== Building server..."

if [ $GOOS = windows ]; then
    go build -o build/server/server-${GOOS}-${GOARCH}.exe cmd/server/main.go
else
    go build -o build/server/server-${GOOS}-${GOARCH} cmd/server/main.go
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
