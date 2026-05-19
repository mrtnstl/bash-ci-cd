#!/usr/bin/env bash

source "$PWD/config.sh"

CURRENT_WORKFLOW_LOG_NAME=$(ls -t "$TEMP_DIR" | head -1)

cd "$WORKDIR/repo"

LOG_FILE_NAME="02_install_deps_log.txt"
#echo "Code version: $(git rev-parse --short HEAD)"
# --silent

if [[ -f package.json && -f package-lock.json ]]; then
    echo "[$(date --utc +%FT%TZ)]___install_deps.sh___" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt" && npm ci &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"
else
    echo "[$(date --utc +%FT%TZ)]___install_deps.sh___package.json and/or package-lock.json not found!" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"
fi
