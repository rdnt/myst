#!/bin/sh
go build -o builder ./cmd/builder
./builder bootstrap
