#//bin/bash

APP_VERSION=$(cat VERSION)
GIT_HASH=$(git rev-parse --short HEAD)
BUILD_TIME="`date +%Y/%m/%d-%H:%M:%S%z`"

echo "APP_VERSION: $APP_VERSION"
echo "GIT_HASH: $GIT_HASH"
echo "BUILD_TIME: $BUILD_TIME"

GO_VERSION_SANITIZED=$(echo "$GO_VERSION" | awk '{print $3 "_" $4}')

ldflags="-X main.AppVersion=$APP_VERSION -X main.GitHash=$GIT_HASH -X main.BuildTime=$BUILD_TIME"
CGO_ENABLED=0 go build -o ./dist/unirest.exe -x --ldflags="$ldflags" && echo 'Build success'
