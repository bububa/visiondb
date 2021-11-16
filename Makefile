PROJECT_NAME=visiondb
GIT_TAG = $(shell git tag | grep ^v | sort -V | tail -n 1)
GIT_REVISION = $(shell git rev-parse --short HEAD)
GIT_SUMMARY = $(shell git describe --tags --dirty --always)
GO_IMPORT_PATH=github.com/bububa/visiondb/server/app
LDFLAGS = -X $(GO_IMPORT_PATH).GitTag=$(GIT_TAG) -X $(GO_IMPORT_PATH).GitRevision=$(GIT_REVISION) -X $(GO_IMPORT_PATH).GitSummary=$(GIT_SUMMARY)

all: *

proto: 
	protoc --go_out=./pb ./pb/*.proto	
server:
	go build -ldflags "$(LDFLAGS)" -o $(PROJECT_NAME)

server-vulkan:
	go build -tags=vulkan -ldflags "$(LDFLAGS)" -o $(PROJECT_NAME)

