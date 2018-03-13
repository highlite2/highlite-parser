FROM golang:1.10.0

RUN curl https://glide.sh/get | sh && \
    go get github.com/golang/lint/golint && \
    go get golang.org/x/tools/cmd/goimports && \
    go get github.com/vektra/mockery/.../

WORKDIR /go/src/highlite2-import