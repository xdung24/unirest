#//bin/bash

GO_OS=$(go env GOOS)
GO_ARCH=$(go env GOARCH)
GO_VERSION=$(go version)
APP_VERSION=$(cat VERSION)
GIT_HASH=$(git rev-parse --short HEAD)
BUILD_TIME="`date +%Y/%m/%d-%H:%M:%S%z`"
GORILLA_VERSION=$(go list -m all | grep github.com/gorilla/mux | awk '{print $2}')

echo "GO_OS: $GO_OS"
echo "GO_ARCH: $GO_ARCH"
echo "GO_VERSION: $GO_VERSION"
echo "APP_VERSION: $APP_VERSION"
echo "GIT_HASH: $GIT_HASH"
echo "BUILD_TIME: $BUILD_TIME"
echo "GORILLA_VERSION: $GORILLA_VERSION"

GO_VERSION_SANITIZED=$(echo "$GO_VERSION" | awk '{print $3 "_" $4}')

ldflags="-X main.GoOs=$GO_OS -X main.GoArch=$GO_ARCH -X main.GoVersion=$GO_VERSION_SANITIZED -X main.AppVersion=$APP_VERSION -X main.GitHash=$GIT_HASH -X main.BuildTime=$BUILDTIME -X main.GinVersion=$GIN_VERSION -X main.GorillaVersion=$GORILLA_VERSION -X main.EchoVersion=$ECHO_VERSION -X main.FiberVersion=$FIBER_VERSION"
go build -o ./tmp/main.exe -x --ldflags="$ldflags"