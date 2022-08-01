#!/usr/bin/env bash

if [[ $EUID != 0 ]]; then
    echo 'You need superuser privileges for installation'
    exit 1
fi

find ./src/ -name "__*" -exec cp {} /bin/ \; 2>/dev/null
