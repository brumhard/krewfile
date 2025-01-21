VERSION 0.8

FROM golang:1.22
ARG --global BINPATH=/usr/local/bin/
ARG --global GOCACHE=/go-cache

ARG --global GOLANGCI_LINT_VERSION=v1.61.0

deps:
    WORKDIR /src
    ENV GO111MODULE=on
    ENV CGO_ENABLED=0
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
    FROM +deps
    COPY main.go .
    ARG GOOS=linux
    ARG GOARCH=amd64
    RUN --mount=type=cache,target=$GOCACHE \
        go build -ldflags="-w -s" -o krewfile main.go
    SAVE ARTIFACT krewfile

lint:
    FROM +deps
    COPY +golangci-lint/golangci-lint $BINPATH
    COPY *.go .
    ARG GOLANGCI_LINT_CACHE=/golangci-cache
    RUN --mount=type=cache,target=$GOCACHE \
        --mount=type=cache,target=$GOLANGCI_LINT_CACHE \
        golangci-lint run -v ./...

test:
    FROM +deps
    COPY *.go .
    ARG GO_TEST="go test"
    RUN --mount=type=cache,target=$GOCACHE \
        $GO_TEST ./...

e2e:
    COPY +krew/krew $BINPATH
    COPY +build/krewfile $BINPATH
    RUN echo "stern" > /root/.krewfile
    RUN krewfile
    RUN krew list 2>/dev/null | grep "stern" >/dev/null

###########
# helper
###########

golangci-lint:
    FROM golangci/golangci-lint:$GOLANGCI_LINT_VERSION
    SAVE ARTIFACT /usr/bin/golangci-lint

krew:
    FROM +deps
    RUN go install sigs.k8s.io/krew/cmd/krew@v0.4.4
    SAVE ARTIFACT /go/bin/krew