.DEFAULT_GOAL := help

## help: Print help information
.PHONYU: help
help:
	@echo "usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | sort | column -t -s ':' |  sed -e 's/^/ /' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## setup: Initialize warehouse setup git scripts and other initialization scripts
.PHONY: setup
setup:
	@sh ./script/setup.sh

## check: Check code formatting and introduce optimizations
.PHONY: check
check:
	@$(MAKE) --no-print-directory tidy
	@$(MAKE) --no-print-directory fmt

## tidy: go mod tidy command...
.PHONY: tidy
tidy:
	@go mod tidy -v

## fmt: go fmt code style...
.PHONY: fmt
fmt:
	@sh ./script/fmt.sh

## api-wire: generate api wire code...
.PHONY: api-wire
api-wire:
	cd $(shell pwd)/internal/api && wire