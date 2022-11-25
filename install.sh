#!/usr/bin/env bash

if [[ -d ./configfiles/ ]]; then
    find ./configfiles/ -name ".*" -exec cp {} $HOME/ \; 2>/dev/null
fi

if [[ $EUID != 0 ]]; then
    echo 'You need superuser privileges to install binaries'
    exit 1
fi

find ./src/ -name "__*" -exec cp {} /bin/ \; 2>/dev/null
