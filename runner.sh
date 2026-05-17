#!/usr/bin/env bash

source ./common/colors.sh
source ./common/spinner.sh

run_with_spinner() {
    local name="$1"
    local command="$2"

    eval "$command" &
    local pid=$! # pid of the last background task

    spinner $pid "$name"

    # wait for task to end
    wait $pid
    local status=$? # exit code

    if [ $status -eq 0 ]; then
        echo -e "$(set_color $(set_color "success" "$BOLD_INTNS_WHITE") "$INTNS_BG_GREEN") ${name} (code: $status)"
    else
        IFS='%'
        echo -e "$(set_color $(set_color "fail   " "$BOLD_INTNS_WHITE") "$INTNS_BG_RED") ${name} (code: $status)"
        unset IFS
    fi

    return $status
}
