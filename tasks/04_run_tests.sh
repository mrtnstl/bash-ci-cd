#!/usr/bin/env bash

source "$PWD/config.sh"

cd "$WORKDIR/repo"
npm install --silent --save-dev jest-silent-reporter

# NODE_OPTIONS flags needed for node.js 25+, otherwise it writes warnings to stdout 
# --no-experimental-webstorage or --no-warnings
NODE_OPTIONS="--no-warnings" npx jest --silent \
  --coverage \
  --coverageDirectory="$WORKDIR/artefacts/test-coverage" \
  --coverageReporters json \
  --coverageReporters html \
  --coverageReporters lcov \
  --reporters jest-silent-reporter \
  --noStackTrace
