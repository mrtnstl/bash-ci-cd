#!/usr/bin/env bash

source "$PWD/config.sh"

CURRENT_WORKFLOW_LOG_NAME=$(ls -t "$TEMP_DIR" | head -1)

cd "$WORKDIR/repo"
# npm install --silent --save-dev jest-silent-reporter

# NODE_OPTIONS flags needed for node.js 25+, otherwise it writes warnings to stdout 
# --no-experimental-webstorage or --no-warnings
# --silent --noStackTrace NODE_OPTIONS="--no-experimental-webstorage" 
echo "[$(date --utc +%FT%TZ)]___run_tests.sh___" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt" && npx jest --coverage \
  --coverageDirectory="$ARTIFACTS_DIR/test-coverage" \
  --coverageReporters json \
  --coverageReporters html \
  --coverageReporters lcov \
  --reporters default &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"
