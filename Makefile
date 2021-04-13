SRC = $(wildcard *.go)

BUILD = build
BUILD_DIRS = ${BUILD}

FIX_SRC = $(foreach val, ${SRC}, ${BUILD}/fix-$(val))

TEST = ${BUILD}/test

.PHONY: fix
fix: ${BUILD_DIRS} ${FIX_SRC} ## Will throw an error if `go fix` finds anything
	@test $(shell wc -l build/fix-*.go | awk '{print $$1}' | paste -sd+ - | bc) = 0

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

${FIX_SRC}: ${BUILD}/fix-%: %
	@cat $< | go fix > $@

${BUILD_DIRS}:
	@mkdir -p $@

.PHONY: test
test: ${BUILD_DIRS} ${TEST} ## Runs go test on all source

${TEST}: ${SRC}
	@go test -v ./... | tee $@

.PHONY: lint
lint: ## Runs goint on all packages, returns an error if it finds anything
	@golint -set_exit_status ./...

.PHONY: clean
clean: ## Cleans all build directories
	@echo removing ${BUILD_DIRS}
	@rm -rf ${BUILD_DIRS}
