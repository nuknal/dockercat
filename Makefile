.PHONY: build clean deploy

BUILD_TIME=`date '+%Y%m%d%H%M%S'`
BUILD_VERSION=${SHANLAAN_MALL_VERSION}
COMMIT_ID=`git rev-parse HEAD`
GO_VERSION=`go version`
BUILD_NAME=version dev
VERSION_PKG=nuknal/dockercat/version
LD_FLAGS="-s -w -X '${VERSION_PKG}.BuildTime=${BUILD_TIME}'                \
                         -X '${VERSION_PKG}.BuildVersion=${BUILD_VERSION}' \
                         -X '${VERSION_PKG}.BuildName=${BUILD_NAME}'       \
                         -X '${VERSION_PKG}.CommitID=${COMMIT_ID}'         \
                         -X '${VERSION_PKG}.GoVersion=${GO_VERSION}'"

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags=${LD_FLAGS} -o bin/dockercat-linux-amd64 main.go
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/dockercat main.go
clean:
	rm -rf ./bin