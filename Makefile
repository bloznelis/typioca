TARGET_SYSTEM ?= $(OS)

ifneq (,$(filter Windows%,$(TARGET_SYSTEM)))
  EXT =.exe
else
  EXT =
endif

GO_FLAGS   ?=
NAME       := typioca
OUTPUT_BIN ?= execs/$(NAME)$(ARCH)$(EXT)
PACKAGE    := github.com/bloznelis/$(NAME)
GIT_REV     = $(shell git rev-parse --short HEAD)
VERSION     = $(shell git describe --abbrev=0 --tags)

default: help

build-win:  ## Builds the win-amd64 CLI
	@env GOOS=windows GOARCH=amd64 ARCH=-win-amd64 make build

build-mac-amd:  ## Builds the mac-amd64 CLI
	@env GOOS=darwin GOARCH=amd64 ARCH=-mac-amd64 make build

build-mac-arm:  ## Builds the mac-arm64 CLI
	@env GOOS=darwin GOARCH=arm64 ARCH=-mac-arm64 make build

build-linux-amd:  ## Builds the linux-amd64 CLI
	@env GOOS=linux GOARCH=amd64 ARCH=-linux-amd64 make build

build:  ## Builds the CLI
	@go build -trimpath ${GO_FLAGS} \
	-ldflags "-w -s -X 'github.com/bloznelis/typioca/cmd.Version=${VERSION}'" \
	-a -tags netgo -o ${OUTPUT_BIN}

build-all: build-win build-mac-amd build-mac-arm build-linux-amd ## Builds execs for all architectures

help: ## This message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'

