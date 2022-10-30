SHELL=/bin/bash

.PHONY: build deploy clean

PATH_CURRENT := $(shell pwd)
PATH_BUILT := $(PATH_CURRENT)/build/server

all:
	deploy

build:
	go build -v -o ./build/server

deploy:
	@if [ "$(wildcard $(PATH_BUILT))" != "" ]; then\
		gcloud run deploy --source . --region asia-southeast1; \
		echo "Done deploy." ; \
		rm -fr "${PATH_BUILT}"; \
		echo "Clean built."; \
	else \
		echo "Please build first."; \
	fi

clean:
	@if [ "$(wildcard $(PATH_BUILT))" != "" ]; then\
		rm -fr "${PATH_BUILT}"; \
		echo "Clean built."; \
	else \
		echo "File does not exist."; \
	fi
