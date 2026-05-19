#!/usr/bin/env bash

source "$PWD/config.sh"

CURRENT_WORKFLOW_LOG_NAME=$(ls -t "$TEMP_DIR" | head -1)

cd "$WORKDIR/repo"
# --silent
echo "[$(date --utc +%FT%TZ)]___test_build.sh___" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt" && npm run build  &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"

