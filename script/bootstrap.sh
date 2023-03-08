#!/bin/sh
go build -o myst ./cmd/builder
./myst bootstrap
