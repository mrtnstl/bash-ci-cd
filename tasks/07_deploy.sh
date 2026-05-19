#!/usr/bin/env bash

source "$PWD/config.sh"

CURRENT_WORKFLOW_LOG_NAME=$(ls -t "$TEMP_DIR" | head -1)

echo "[$(date --utc +%FT%TZ)]___deploy.sh___NOT_IMPLEMENTED" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"