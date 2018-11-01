NAME=rubixcore
REPO=github.com/hackstock/${NAME}

BINARY=${NAME}
BINARY_SRC=$(REPO)/cmd/${NAME}
BUILD_DIR?=$(CURDIR)/out
GOOS ?= linux
GOARCH ?= amd64
GO_LINKER_FLAGS=-ldflags="-s -w"

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

.PHONY:  build

build:
	@mkdir -p ${BUILD_DIR}
	@printf "${OK_COLOR}==> Building binary into ${BUILD_DIR}${NO_COLOR}\n"
	@CGO_ENABLED=0 go build -o ${BUILD_DIR}/${BINARY} ${GO_LINKER_FLAGS} ${BINARY_SRC}

test-unit:
	@printf "${OK_COLOR}==> Running unit tests${NO_COLOR}\n"
	@go test -count=1 -v -race -coverprofile=coverage.txt --covermode=atomic ./...

clean:
	@printf "${OK_COLOR}==> Cleaning project${NO_COLOR}\n"
	if [ -d ${BUILD_DIR} ] ; then rm -rf ${BUILD_DIR}/* ; fi
