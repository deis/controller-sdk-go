# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/controller-sdk-go

REVISION ?= $(shell git rev-parse --short HEAD)
REGISTRY ?= quay.io/
IMAGE_PREFIX ?= deisci
IMAGE := ${REGISTRY}${IMAGE_PREFIX}/controller-sdk-go-dev:${REVISION}

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.16.0
DEV_ENV_WORK_DIR := /go/src/${repo_path}
DEV_ENV_PREFIX := docker run --rm -e CGO_ENABLED=0 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_PREFIX_CGO_ENABLED := docker run --rm  -e CGO_ENABLED=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

# UID and GID of local user
UID := $(shell id -u)
GID := $(shell id -g)

GOTEST = go test --cover --race

bootstrap:
	${DEV_ENV_CMD} glide install

glideup:
	${DEV_ENV_CMD} glide up

test-style: build-test-image
	docker run --rm ${IMAGE} lint

test: build-test-image
	docker run --rm ${IMAGE} test

build-test-image:
	docker build -t ${IMAGE} .

push-test-image: build-test-image
