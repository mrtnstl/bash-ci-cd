#!/usr/bin/env bash

run_with_spinner() {
    local name="$1"
    local command="$2"

    eval "$command" &
    local pid=$! # pid of the last background task

    spinner $pid "$name"

    # wait for task to end
    wait $pid
    local status=$? # exit code

    IFS=''

    if [ $status -eq 0 ]; then
        NOTIFICATION_BODY="${NOTIFICATION_BODY}<tr style='line-height: 0px;'><td><p style='margin:0; color:#9a9a9a;'>${name}</p></td><td><p>&#9989;</p></td></tr>"
        echo -e "$(set_color $(set_color "    ok " "$BOLD_INTNS_WHITE") "$INTNS_BG_GREEN") ${name} (code: $status)"
    else
        NOTIFICATION_BODY="${NOTIFICATION_BODY}<tr style='line-height: 0px;'><td><p style='margin:0; color:#9a9a9a;'>${name}</p></td><td><p>&#10060;</p></td></tr>"
        echo -e "$(set_color $(set_color "  fail " "$BOLD_INTNS_WHITE") "$INTNS_BG_RED") ${name} (code: $status)"
    fi

    unset IFS

    return $status
}
