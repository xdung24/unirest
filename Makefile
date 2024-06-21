GO_OS := $(shell go env GOOS)
GO_ARCH := $(shell go env GOARCH)

build:
	GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -o universal-rest ./...

clean:
	rm -f universal-rest

.PHONY: build clean