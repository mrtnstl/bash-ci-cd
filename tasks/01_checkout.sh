#!/usr/bin/env bash

source "$PWD/config.sh"

CURRENT_WORKFLOW_LOG_NAME=$(ls -t "$TEMP_DIR" | head -1)

cd "$WORKDIR"

if [ ! -d "repo" ]; then
    git clone -q -b "$BRANCH" "$REPO_URL" repo
    echo "[$(date --utc +%FT%TZ)]___checkout.sh completed___repo cloned" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"
else 
    cd repo
    git fetch -q origin
    git reset -q --hard origin/"$BRANCH"
    echo "[$(date --utc +%FT%TZ)]___checkout.sh completed___repo fetched and reset" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"
fi