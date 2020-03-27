#!/bin/bash
echo "Bootstrapping Myst..."
sleep 1
echo "Installing NPM modules..."
cd vue > /dev/null 2>&1
npm install > /dev/null 2>&1
cd ../ > /dev/null 2>&1
echo "Bootstrapping complete."
sleep 3
