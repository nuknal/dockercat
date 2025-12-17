.PHONY: build clean deploy

BUILD_TIME=`date '+%Y%m%d%H%M%S'`
BUILD_VERSION=0.1.1
COMMIT_ID=`git rev-parse HEAD`
GO_VERSION=`go version`
BUILD_NAME=release
VERSION_PKG=github.com/nuknal/dockercat/pkg/version
LD_FLAGS="-s -w -X '${VERSION_PKG}.BuildTime=${BUILD_TIME}'                \
                         -X '${VERSION_PKG}.BuildVersion=${BUILD_VERSION}' \
                         -X '${VERSION_PKG}.BuildName=${BUILD_NAME}'       \
                         -X '${VERSION_PKG}.CommitID=${COMMIT_ID}'         \
                         -X '${VERSION_PKG}.GoVersion=${GO_VERSION}'"

build: 
	env GOOS=darwin GOARCH=amd64 go build -ldflags=${LD_FLAGS} -o bin/dockercat main.go

clean:
	rm -rf ./bin