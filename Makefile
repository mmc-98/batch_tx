GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION:=$(shell git describe --tags --always)



ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif


.PHONY: build.generate
# generate generate
build.generate:
		go build   -ldflags="-s -w"  -ldflags "-X main.Version=$(VERSION)"  -ldflags "-X main.Name=generate" -o build/generate cmd/generate/*.go


.PHONY: build.send
# generate sendo
build.send:
		go build   -ldflags="-s -w"  -ldflags "-X main.Version=$(VERSION)"  -ldflags "-X main.Name=send" -o build/send cmd/send/*.go



.PHONY: build.batch_tx
# generate batch_tx
build.batch_tx:
		go build   -ldflags="-s -w"  -ldflags "-X main.Version=$(VERSION)"  -ldflags "-X main.Name=batch_tx" -o build/batch_tx *.go

.PHONY: build.config
# generate config
build.config:
		mkdir -p build
		cp -rf etc  build


.PHONY: all
# generate all
all:
	go mod tidy
	make build.config
	make build.generate;
	make build.send;
	make build.batch_tx;


.PHONY: start.generate
# start.generate
start.generate:
		./build/generate

.PHONY: start.send
# start.send
start.send:
		./build/send

.PHONY: start.tx
# start.tx
start.tx:
		./build/batch_tx

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, sdfsdfsdfdsfsdfsfRSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
