#!/usr/bin/env bash

source "$PWD/config.sh"

# deleting artifact and temp dir
rm -rf "$ARTIFACTS_DIR"
rm -rf "$TEMP_DIR"

# creating directories for workflow
mkdir -p "$WORKDIR"
mkdir -p "$LOGS_DIR"
mkdir -p "$ARTIFACTS_DIR"
mkdir -p "$TEMP_DIR"

# creating a temp file with the current datetime for log file naming
touch "$TEMP_DIR/$(date --utc +%FT%TZ)_workflow_log"

CURRENT_WORKFLOW_LOG_NAME=$(ls -t "$TEMP_DIR" | head -1)

# echo "$CURRENT_WORKFLOW_LOG_NAME"

echo "[$(date --utc +%FT%TZ)]___init.sh completed" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"