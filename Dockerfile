FROM quay.io/deis/go-dev:0.17.0
# This Dockerfile is used to bundle the source and all dependencies into an image for testing.

COPY glide.yaml /go/src/github.com/deis/controller-sdk-go/
COPY glide.lock /go/src/github.com/deis/controller-sdk-go/

WORKDIR /go/src/github.com/deis/controller-sdk-go

RUN glide install --strip-vcs --strip-vendor

COPY ./_scripts /usr/local/bin

COPY . /go/src/github.com/deis/controller-sdk-go
