#!/usr/bin/env bash

source "$PWD/config.sh"

cd "$WORKDIR/repo"
npm install --silent --save-dev jest-silent-reporter
#npx jest --silent --coverage --coverageReporters="json,html,lcov" --reporters=jest-silent-reporter

npx jest --silent \
  --coverage \
  --coverageReporters="json" \
  --coverageReporters="html" \
  --coverageReporters="lcov" \
  --reporters=jest-silent-reporter \
  --noStackTrace

#npx jest --silent # --coverage --coverageDirectory="$WORKDIR/artefacts/test-coverage" --silent
#mv "$WORKDIR/repo/coverage" "$WORKDIR/artefacts"