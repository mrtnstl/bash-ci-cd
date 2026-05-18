#!/usr/bin/env bash

source "$PWD/config.sh"

cd "$WORKDIR/repo"

#echo "Code version: $(git rev-parse --short HEAD)"

if [[ -f package.json && -f package-lock.json ]]; then
    npm ci --silent
else
    echo "package.json and=or package-lock.json not found!"
fi