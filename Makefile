export GO15VENDOREXPERIMENT=1

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/controller-sdk-go

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.10.0
DEV_ENV_WORK_DIR := /go/src/${repo_path}
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -e CGO_ENABLED=0 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_PREFIX_CGO_ENABLED := docker run --rm -e GO15VENDOREXPERIMENT=1 -e CGO_ENABLED=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

GO_FILES = $(wildcard *.go)
GO_PACKAGES = api apps auth builds certs config domains keys perms ps releases users pkg/time
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))
GOFMT = gofmtresult=$$(gofmt -e -l -s ${GO_FILES} ${GO_PACKAGES}); if [[ -n $$gofmtresult ]]; then echo "gofmt errors found in the following files: $${gofmtresult}"; false; fi;
GOTEST = go test --cover --race -v

bootstrap:
	${DEV_ENV_CMD} glide install

glideup:
	${DEV_ENV_CMD} glide up

setup-gotools:
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/cover
	go get -u golang.org/x/tools/cmd/vet

test: test-style test-unit

test-style:
	${DEV_ENV_CMD} bash -c '${GOFMT}'
	${DEV_ENV_CMD} sh -c 'go vet $(repo_path) $(GO_PACKAGES_REPO_PATH)'
	@for i in $(addsuffix /...,$(GO_PACKAGES)); do \
		${DEV_ENV_CMD} golint $$i; \
	done

test-unit:
	${DEV_ENV_PREFIX_CGO_ENABLED} ${DEV_ENV_IMAGE} sh -c '${GOTEST} $$(glide nv)'
