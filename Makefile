# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/controller-sdk-go

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

test: test-style test-unit

test-style:
	${DEV_ENV_CMD} lint

test-unit:
	${DEV_ENV_PREFIX_CGO_ENABLED} ${DEV_ENV_IMAGE} sh -c '${GOTEST} $$(glide nv)'

test-cover: test-style
	${DEV_ENV_PREFIX_CGO_ENABLED} ${DEV_ENV_IMAGE} test-cover.sh

# Set local user as owner for files
fileperms:
	${DEV_ENV_PREFIX_CGO_ENABLED} ${DEV_ENV_IMAGE} chown -R ${UID}:${GID} .
