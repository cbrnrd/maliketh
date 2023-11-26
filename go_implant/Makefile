GO=go
GOFLAGS=
DEBUG ?= 1
BINARY=implant
MAIN=cmd/main.go
GOARCH ?= amd64

# If debug, append -debug to binary name
ifeq ($(DEBUG),1)
	BINARY:=$(BINARY)-debug
else
	BINARY:=$(BINARY)-release
endif

ifeq ($(DEBUG),1)
	GOFLAGS:=$(GOFLAGS) -gcflags="all=-N -l" -ldflags="-X maliketh.config.DEBUG=true -X maliketh.config.INITIAL_SLEEP=1"
else
	GOFLAGS:=$(GOFLAGS) -ldflags="-s -w -X maliketh.config.DEBUG=false" -trimpath -gcflags=all="-l -B"
endif


default: native-debug

native-debug:
	$(GO) build $(GOFLAGS) -o implant-dev-debug $(MAIN)

build:
	$(GO) build $(GOFLAGS) -o implant-$(shell go env GOOS)-$(GOARCH) $(MAIN)

all: setup deps macos linux windows

setup:
	mkdir -p build

deps:
	$(GO) mod tidy

linux:
	@/bin/echo -n "-----> Building Linux binary ($(GOARCH))... "
	@GOOS=linux GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o build/$(BINARY)-linux-$(GOARCH) $(MAIN)
	@echo "DONE"

macos:
	@/bin/echo -n "-----> Building MacOS binary ($(GOARCH))... "
	@GOOS=darwin GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o build/$(BINARY)-macos-$(GOARCH) $(MAIN)
	@echo "DONE"

windows:
	@/bin/echo -n "-----> Building Windows binary ($(GOARCH))... "
	@GOOS=windows GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o build/$(BINARY)-windows-$(GOARCH).exe $(MAIN)
	@echo "DONE"

clean:
	rm -rf build/*