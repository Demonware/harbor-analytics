# Copyright (C) Activision Publishing, Inc. 2017
# https://github.com/Demonware/harbor-analytics
# Author: David Rieger
# Licensed under the 3-Clause BSD License (the "License");
# you may not use this file except in compliance with the License.

clean:
	rm -rf ../out

get-deps:
	go get -d
	#go get github.com/jung-kurt/gofpdf
	#go get github.com/wcharczuk/go-chart

update-self:
	-rm -rf $(shell echo $$GOPATH | cut -d: -f 1)/src/github.com/Demonware/harbor-analytics/analyst
	mkdir -p $(shell echo $$GOPATH | cut -d: -f 1)/src/github.com/Demonware/harbor-analytics/analyst
	cp -rf . $(shell echo $$GOPATH | cut -d: -f 1)/src/github.com/Demonware/harbor-analytics/analyst

build: clean update-self get-deps
	go build
