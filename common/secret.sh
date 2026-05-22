#!/usr/bin/env bash

get_secret(){
    local key=$1 value=$2
    result=$(secret-tool lookup $key $value)

    echo "$result"
}

set_secret(){
    local label=$1 key=$2 value=$3 secret=$4

    echo "$secret" | secret-tool store --label="$1" $key $value
    local result=$?
    
    if [ $result -ne 0 ]; then
        echo "error: secret-tool exited with non zero value"
    fi

}
