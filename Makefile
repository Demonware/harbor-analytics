# Copyright (C) Activision Publishing, Inc. 2017
# https://github.com/Demonware/harbor-analytics
# Author: David Rieger
# Licensed under the 3-Clause BSD License (the "License");
# you may not use this file except in compliance with the License.

.PHONY: check-pre-get-raw-data
check-pre-get-raw-data:
ifndef AWS_ACCESS_KEY_ID
		$(error AWS_ACCESS_KEY_ID is undefined)
endif
ifndef AWS_SECRET_ACCESS_KEY
		$(error AWS_SECRET_ACCESS_KEY is undefined)
endif

.PHONY: build-raw-data-fetcher
build-raw-data-fetcher:
	docker build -t harboranalyst/data-fetcher docker-get-raw-data

.PHONY: get-raw-data
get-raw-data: check-pre-get-raw-data
	-rm -rf ./raw
	mkdir raw

	docker run --rm -e 'AWS_ACCESS_KEY_ID=$(shell echo $$AWS_ACCESS_KEY_ID)' -e 'AWS_SECRET_ACCESS_KEY=$(shell echo $$AWS_SECRET_ACCESS_KEY)' -e "AWS_DEFAULT_REGION=us-west-2" -v "$(PWD)/raw:/raw" harboranalyst/data-fetcher

.PHONY: build
build:
	docker build -t harboranalyst/analyst -f docker-analyst/Dockerfile .

.PHONY: run
run:
	-rm -rf ./out
	mkdir out
	docker run --rm -v $(PWD)/analyst.yaml:/root/analyst.yaml -v $(PWD)/raw:/root/raw -v $(PWD)/out:/root/out harboranalyst/analyst

.PHONY: check-pre-publish
check-pre-publish:
ifndef SLACK_TOKEN
		$(error SLACK_TOKEN is undefined)
endif

.PHONY: publish
publish: check-pre-publish
	cd slack && ./send_report.sh $(SLACK_TOKEN)
