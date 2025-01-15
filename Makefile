BIN_OUTPUT_PATH = bin
TOOL_BIN = bin/gotools/$(shell uname -s)-$(shell uname -m)
UNAME_S ?= $(shell uname -s)
GOPATH = $(HOME)/go/bin
export PATH := ${PATH}:$(GOPATH)
GOIMPORTS ?= $(shell which goimports)

build: format update-rdk
	rm -f $(BIN_OUTPUT_PATH)/transform-camera-extended
	go build $(LDFLAGS) -o $(BIN_OUTPUT_PATH)/transform-camera-extended main.go

module.tar.gz: build
	rm -f $(BIN_OUTPUT_PATH)/module.tar.gz
	tar czf $(BIN_OUTPUT_PATH)/module.tar.gz $(BIN_OUTPUT_PATH)/transform-camera-extended meta.json

setup:
	if [ "$(UNAME_S)" = "Linux" ]; then \
		sudo apt-get install -y apt-utils coreutils tar libnlopt-dev libjpeg-dev pkg-config; \
	fi
	# remove unused imports
	# go install golang.org/x/tools/cmd/goimports@latest
	# find . -name '*.go' -exec $(GOIMPORTS) -w {} +


clean:
	rm -rf $(BIN_OUTPUT_PATH)/transform-camera-extended $(BIN_OUTPUT_PATH)/module.tar.gz transform-camera-extended

format:
	gofmt -w -s .

update-rdk:
	go get go.viam.com/rdk@latest
	go mod tidy
