#!/usr/bin/env bash

# Copyright (C) Activision Publishing, Inc. 2017
# https://github.com/Demonware/harbor-analyst
# Author: David Rieger
# Licensed under the 3-Clause BSD License (the "License");
# you may not use this file except in compliance with the License.

TOKEN=$1

if [ -z $TOKEN ]
then
    echo "Please provide a slack token as first argument. Abort." >&2
    exit 1
fi

curl -F file=@../out/report.pdf -F channels="#channel" -F token=$TOKEN https://slack.com/api/files.upload
