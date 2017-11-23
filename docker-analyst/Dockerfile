# Copyright (C) Activision Publishing, Inc. 2017
# https://github.com/Demonware/harbor-analyst
# Author: David Rieger
# Licensed under the 3-Clause BSD License (the "License");
# you may not use this file except in compliance with the License.

FROM golang:1.9-alpine3.6
RUN apk add --update make git
ADD analyst /analyst
RUN cd /analyst && GOOS=linux make build

FROM alpine:3.6
WORKDIR /root
RUN mkdir /root/bin
COPY --from=0 /analyst/analyst /root/bin/analyst
RUN chmod +x /root/bin/analyst
WORKDIR /root/bin
ENTRYPOINT ["./analyst"]
