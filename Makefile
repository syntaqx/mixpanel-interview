# Use bash as shell (Note: Ubuntu now uses dash which doesn't support PIPESTATUS)
SHELL := $(shell which bash)

# defaults to using -s, unless VERBOSE is set
ifeq ($(VERBOSE)_x, _x)
	MAKEFLAGS+=-s
endif

# delete built-in suffixes and define .go
.SUFFIXES:
.SUFFIXES: .go

# Set the default target to the help text
.DEFAULT: help

# Ensure targets never collide with a file.
.PHONY: help dep clean build run test test-cover bench

help: ## Show help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dep: ## Install build dependencies
	go get -v -u github.com/golang/dep/cmd/dep
	dep ensure

clean: ## Clean build artifacts
	rm -rf bin *.out

build: ## Build the application
	CGO_ENABLED=0 go build -o bin/mixpanel cmd/mixpanel/*.go

run: build ## Run the applicaiton
	./bin/mixpanel --debug

test: ## Run unit tests
	go test -v ./...

test-cover: ## Generate test coverage profile
	echo "" > coverage.out
	for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.out && rm profile.out; \
	done

bench: ## Run benchmarks
	go test -bench .
