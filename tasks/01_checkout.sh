#!/usr/bin/env bash

source "$PWD/config.sh"

mkdir -p "$WORKDIR"
cd "$WORKDIR"

if [ ! -d "repo" ]; then
    git clone -q -b "$BRANCH" "$REPO_URL" repo
else 
    cd repo
    git fetch -q origin
    git reset -q --hard origin/"$BRANCH"
fi