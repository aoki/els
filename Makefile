NAME			:= els
VERSION		:= v0.1.0
REVISION	:= $(shell git rev-parse --short HEAD)

SRCS		:= $(shell find src -type f -name '*.go')
LDFLAGS	:= -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

.PHONY: dep-install
dep-install:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/...
endif

.PHONY: init
init:
	dep init

.PHONY: deps
deps: dep
	dep ensure -update

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) $(SRCS)

.PHONY: install
install:
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o ${GOPATH}/bin/$(NAME)
	#go install $(LDFLAGS) -o hello

.PHONY: imports
imports:
	# if you want to install goimports
	# `go get golang.org/x/tools/cmd/goimports`
ifeq ($(shell command -v goimports 2> /dev/null),)
	go get -u github.com/golang/dep/...
endif
	goimports -w *.go
