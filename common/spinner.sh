#!/usr/bin/env bash

spinner() {
    local pid=$1
    local delay=0.1
    local spin='-\|/'
    local i=0

    # while process is running
    while kill -0 $pid 2>/dev/null; do
        printf " ${spin:i++%4:1}  %s" "$2"
        sleep $delay
        printf "                                              \r"
    done
    printf "                                              \r"
}