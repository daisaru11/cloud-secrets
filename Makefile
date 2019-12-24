export GO111MODULE=on

build:
	go build

test:
	go test ./... -v

lint:
	golangci-lint run

build_image:
	docker build -t daisaru11/cloud-secrets:0.0.1 .

.PHONY: build test lint build_image

