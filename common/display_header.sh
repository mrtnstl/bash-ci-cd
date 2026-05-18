#!/usr/bin/env bash

display_header()
{
    IFS='%'

    HEADER_COLOR="$BOLD_INTNS_RED"
    h_a_length=0

    header_ascii=(
        "$(set_color " ███████████                    █████           █████████  █████        ███    █████████  ██████████     " "$HEADER_COLOR")"
        "$(set_color "░░███░░░░░███                  ░░███           ███░░░░░███░░███        ███░   ███░░░░░███░░███░░░░███    " "$HEADER_COLOR")"
        "$(set_color " ░███    ░███  ██████    █████  ░███████      ███     ░░░  ░███       ███░   ███     ░░░  ░███   ░░███   " "$HEADER_COLOR")"
        "$(set_color " ░██████████  ░░░░░███  ███░░   ░███░░███    ░███          ░███      ███░   ░███          ░███    ░███   " "$HEADER_COLOR")"
        "$(set_color " ░███░░░░░███  ███████ ░░█████  ░███ ░███    ░███          ░███     ███░    ░███          ░███    ░███   " "$HEADER_COLOR")"
        "$(set_color " ░███    ░███ ███░░███  ░░░░███ ░███ ░███    ░░███     ███ ░███    ███░     ░░███     ███ ░███    ███    " "$HEADER_COLOR")"
        "$(set_color " ███████████ ░░████████ ██████  ████ █████    ░░█████████  █████  ███░       ░░█████████  ██████████     " "$HEADER_COLOR")"
        "$(set_color "░░░░░░░░░░░   ░░░░░░░░ ░░░░░░  ░░░░ ░░░░░      ░░░░░░░░░  ░░░░░  ░░░          ░░░░░░░░░  ░░░░░░░░░░      " "$HEADER_COLOR")"
    )

    for row in "${header_ascii[@]}"; do
        echo -e "$row"
        h_a_length=$((h_a_length + 1))
    done

    unset IFS

    echo -e $(set_color "Bash CI/CD (mode: $ENV)" "$BOLD_INTNS_RED")
}