#!/usr/bin/env bash

source ./config.sh
source ./common/colors.sh
source ./common/display_header.sh
source ./runner.sh

display_header

echo ""

# TODO: check if config is valid

echo -e "$(set_color "$(set_color "start  " "$BOLD_INTNS_WHITE")" "$INTNS_BG_BLUE") Starting pipeline"

SECONDS=0

chmod u+x ./tasks/*

for step in tasks/*.sh; do
    run_with_spinner "$(basename $step)" "$step"
    runner_result=$?
    if [ $runner_result -ne 0 ]; then
        printf "\a\a\a"
        duration=$SECONDS
        echo -e "$(set_color "$(set_color "stopped" "$BOLD_INTNS_WHITE")" "$INTNS_BG_RED") Pipeline failed and ran for $((duration / 60)) minutes and $((duration % 60)) seconds."
        exit 1
    else
        printf "\a"
    fi
done

chmod u-x ./tasks/*

duration=$SECONDS

echo -e "$(set_color "$(set_color "stop   " "$BOLD_INTNS_WHITE")" "$INTNS_BG_BLUE") Pipeline ran for $((duration / 60)) minutes and $((duration % 60)) seconds."
