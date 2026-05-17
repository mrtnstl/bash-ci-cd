#!/usr/bin/env bash

# reset code
RESET="\e[0m"

# bold, intense text colors
BOLD_INTNS_BLACK="\e[1;90m"
BOLD_INTNS_RED="\e[1;91m"
BOLD_INTNS_GREEN="\e[1;92m"
BOLD_INTNS_YELLOW="\e[1;93m"
BOLD_INTNS_BLUE="\e[1;94m"
BOLD_INTNS_PURPLE="\e[1;95m"
BOLD_INTNS_CYAN="\e[1;96m"
BOLD_INTNS_WHITE="\e[1;97m"

# intense background colors
INTNS_BG_BLACK="\e[0;100m"
INTNS_BG_RED="\e[0;101m"
INTNS_BG_GREEN="\e[0;102m"
INTNS_BG_YELLOW="\e[0;103m"
INTNS_BG_BLUE="\e[0;104m"
INTNS_BG_PURPLE="\e[0;105m"
INTNS_BG_CYAN="\e[0;106m"
INTNS_BG_WHITE="\e[0;107m"

set_color()
{
    local text=$1 color_code=$2

    echo "$color_code$text$RESET"
}
