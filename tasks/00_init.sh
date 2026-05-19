#!/usr/bin/env bash

# check if config.sh file exists
test ! -e "$PWD/config.sh" \
    && echo "config.sh does not exists" \
    && exit 1

source "$PWD/config.sh"

# check for missing config ($ENV is not important)
if ! echo "$REPO_URL"  >/dev/null 2>&1 || \
   ! echo "$BRANCH" >/dev/null 2>&1 || \
   ! echo "$WORKDIR"   >/dev/null 2>&1 || \
   ! echo "$LOGS_DIR" >/dev/null 2>&1 || \
   ! echo "$ARTIFACTS_DIR"  >/dev/null 2>&1 || \
   ! echo "$TEMP_DIR"  >/dev/null 2>&1; then
    echo "required configuration is missing!"
    exit 1
fi 

# check dependencies
if ! command -v git  >/dev/null 2>&1 || \
   ! command -v curl >/dev/null 2>&1 || \
   ! command -v jq   >/dev/null 2>&1 || \
   ! command -v secret-tool >/dev/null 2>&1; then
    echo "system level dependency missing!"
    exit 1
fi 

# check and set up node if not exists
if ! command -v node >/dev/null 2>&1 || \
   ! command -v npm  >/dev/null 2>&1 || \
   ! command -v npx  >/dev/null 2>&1; then
    echo "node.js is missing!"
    exit 1
fi 

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

echo "[$(date --utc +%FT%TZ)]___init.sh completed" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"