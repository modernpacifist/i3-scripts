#!/usr/bin/env bash

find ./src/ -name "__*" -exec cp {} /bin/ \; 2>/dev/null
