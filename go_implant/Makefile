GO=go
GOFLAGS=
DEBUG ?= 0
BINARY=implant
MAIN=cmd/main.go
GOARCH ?= amd64
GOOS ?= $(shell go env GOOS)


# If debug, append -debug to binary name
ifeq ($(DEBUG),1)
	BINARY:=$(BINARY)-debug
else
	BINARY:=$(BINARY)-release
endif

ifeq ($(DEBUG),1)
	GOFLAGS:=$(GOFLAGS) -gcflags="all=-N -l" -ldflags="-X maliketh.config.DEBUG=true -X maliketh.config.INITIAL_SLEEP=1" -tags debug
else
	GOFLAGS:=$(GOFLAGS) -ldflags="-s -w -X maliketh.config.DEBUG=false" -trimpath -gcflags=all="-l -B"
endif


default: native-debug

native-debug:
	$(GO) build \
	-gcflags="all=-N -l" -ldflags="-X maliketh.config.DEBUG=true -X maliketh.config.INITIAL_SLEEP=1" \
	-tags debug -o implant-dev-debug -tags debug $(MAIN)

build:
	$(GO) build $(GOFLAGS) -o implant-$(shell go env GOOS)-$(GOARCH) $(MAIN)

most: setup deps macos linux windows

setup:
	@mmkdir -p build

deps:
	$(GO) mod tidy

prepare-release:
	@mkdir -p src_release
	@cp -r ./cmd ./pkg ./go.mod ./go.sum src_release
	@find ./src_release -type f -name "*.go" -print0 | xargs -0 sed -i '' 's/DebugPrintln/\/\/DebugPrintln/g'

linux:
	@/bin/echo -n "-----> Building Linux binary ($(GOOS):$(GOARCH), DEBUG=$(DEBUG))... "
	@GOOS=$(or $(GOOS), linux) GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o build/$(BINARY)-$(GOOS)-$(GOARCH) $(MAIN)
	@echo "DONE"

macos:
	@/bin/echo -n "-----> Building MacOS binary ($(GOOS):$(GOARCH), DEBUG=$(DEBUG))... "
	@GOOS=darwin GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o build/$(BINARY)-$(GOOS)-$(GOARCH) $(MAIN)
	@echo "DONE"

windows:
	@/bin/echo -n "-----> Building Windows binary ($(GOOS):$(GOARCH), DEBUG=$(DEBUG))... "
	@GOOS=windows GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o build/$(BINARY)-$(GOOS)-$(GOARCH).exe $(MAIN)
	@echo "DONE"

all: setup deps
# Macos
	@GOOS=darwin GOARCH=arm64 make macos 
	@GOOS=darwin GOARCH=amd64 make macos
		
# # FreeBSD
# 	@GOOS=freebsd GOARCH=386 make linux
# 	@GOOS=freebsd GOARCH=arm64 make linux
# 	@GOOS=freebsd GOARCH=arm make linux
	
# Generic linux
	@GOOS=linux GOARCH=386 make linux
	@GOOS=linux GOARCH=amd64 make linux
	@GOOS=linux GOARCH=arm make linux
	@GOOS=linux GOARCH=arm64 make linux
	@GOOS=linux GOARCH=ppc64 make linux
	@GOOS=linux GOARCH=ppc64le make linux
	@GOOS=linux GOARCH=mips make linux
	@GOOS=linux GOARCH=mipsle make linux
	@GOOS=linux GOARCH=mips64 make linux
	@GOOS=linux GOARCH=mips64le make linux
	
# # NetBSD
# 	@GOOS=netbsd GOARCH=386 make linux
# 	@GOOS=netbsd GOARCH=amd64 make linux
# 	@GOOS=netbsd GOARCH=arm make linux
# 	@GOOS=netbsd GOARCH=386 make linux

# # OpenBSD
# 	@GOOS=openbsd GOARCH=386 make linux
# 	@GOOS=openbsd GOARCH=amd64 make linux
# 	@GOOS=openbsd GOARCH=arm make linux

# Windows
	@GOOS=windows GOARCH=386 make windows
	@GOOS=windows GOARCH=amd64 make windows
	@GOOS=windows GOARCH=arm make windows

release: prepare-release most
	rm -rf ./src_release


clean:
	rm -rf build/*