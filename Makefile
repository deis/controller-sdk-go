# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/controller-sdk-go

REVISION ?= $(shell git rev-parse --short HEAD)
REGISTRY ?= quay.io/
IMAGE_PREFIX ?= deisci
IMAGE := ${REGISTRY}${IMAGE_PREFIX}/controller-sdk-go-dev:${REVISION}

test-style: build-test-image
	docker run --rm ${IMAGE} lint

test-unit: build-test-image
	docker run --rm ${IMAGE} test

test: build-test-image test-style test-unit

build-test-image:
	docker build -t ${IMAGE} .

push-test-image:
	docker push ${IMAGE}
