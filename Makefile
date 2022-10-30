SHELL=/bin/bash

PATH_CURRENT := $(shell pwd)
PATH_BUILT := $(PATH_CURRENT)/build/server

all: build deploy clean

build:
	env GOOS=linux GOARCH=amd64 go build -v -o ./build/server

deploy: build
	gcloud run deploy --source . --region asia-southeast1; \
	echo "Done deploy."

clean:
	rm -fr "${PATH_BUILT}"; \
	echo "Clean built."