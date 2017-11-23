#!/bin/sh

# Copyright (C) Activision Publishing, Inc. 2017
# https://github.com/Demonware/harbor-analytics
# Author: David Rieger
# Licensed under the 3-Clause BSD License (the "License");
# you may not use this file except in compliance with the License.

echo "Start..."
for csv in $(aws s3 ls s3://bucketname/daily/ | grep '[.]csv' | tr -s ' ' | cut -d ' ' -f 4)
do
        aws s3 cp s3://bucketname/daily/$csv /raw
done
