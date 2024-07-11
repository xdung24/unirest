#//bin/bash

GO_OS=$(go env GOOS)
GO_ARCH=$(go env GOARCH)
GO_VERSION=$(go version)
APP_VERSION=$(cat VERSION)
GIT_HASH=$(git rev-parse --short HEAD)
GIN_VERSION=$(go list -m all | grep github.com/gin-gonic/gin | awk '{print $2}')
GORILLA_VERSION=$(go list -m all | grep github.com/gorilla/mux | awk '{print $2}')
ECHO_VERSION=$(go list -m all | grep github.com/labstack/echo/v4 | awk '{print $2}')
FIBER_VERSION=$(go list -m all | grep github.com/gofiber/fiber/v2 | awk '{print $2}')

echo "GO_OS: $GO_OS"
echo "GO_ARCH: $GO_ARCH"
echo "GO_VERSION: $GO_VERSION"
echo "APP_VERSION: $APP_VERSION"
echo "GIT_HASH: $GIT_HASH"
echo "BUILDTIME: $BUILDTIME"
echo "GIN_VERSION: $GIN_VERSION"
echo "GORILLA_VERSION: $GORILLA_VERSION"
echo "ECHO_VERSION: $ECHO_VERSION"
echo "FIBER_VERSION: $FIBER_VERSION"

ldflags="-X 'main.GoOs=$GO_OS' -X 'main.GoArch=$GO_ARCH' -X 'main.GoVersion=$GO_VERSION' -X 'main.AppVersion=$APP_VERSION' -X 'main.GitHash=$GIT_HASH' -X 'main.BuildTime=$BUILDTIME' -X 'main.GinVersion=$GIN_VERSION' -X 'main.GorillaVersion=$GORILLA_VERSION' -X 'main.EchoVersion=$ECHO_VERSION' -X 'main.FiberVersion=$FIBER_VERSION'"
CGO_ENABLED=0 GOOS=$GO_OS GOARCH=$GO_ARCH go build -o ./tmp/main.exe -x --ldflags=$ldflags | true