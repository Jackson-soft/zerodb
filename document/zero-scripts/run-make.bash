#!/usr/bin/env bash

# if any command inside script returns error, exit and return that error 
set -e

cd "${0%/*}/.."

echo "Running make"
echo "............................"
make
