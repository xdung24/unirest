GO_OS := $(shell go env GOOS)
GO_ARCH := $(shell go env GOARCH)

build:
	CGO_ENABLED=0 GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -x -o universal-rest

clean:
	rm -f universal-rest

.PHONY: build clean