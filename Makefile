
# GO_SRCS := $(shell find . -type f -name '*.go' -not -path './examples/*' -not -path './tools/*')
GO_MODULE_FILES := ./go.mod ./go.sum
# DEPS := $(GO_SRCS) $(GO_MODULE_FILES)

# GO_EXAMPLES_SRCS := $(shell find ./examples -type f -name '*.go')
GO_EXAMPLES_MODULE_FILES := ./examples/go.mod ./examples/go.sum
# EXAMPLES_DEPS := $(GO_EXAMPLES_SRCS) $(GO_EXAMPLES_MODULE_FILES)

# GO_TOOLS_SRCS := $(shell find ./tools -type f -name '*.go')
GO_TOOLS_MODULE_FILES := ./tools/go.mod ./tools/go.sum
# TOOLS_DEPS := $(GO_TOOLS_SRCS) $(GO_TOOLS_MODULE_FILES)

.PHONY: all
all: lint test

.PHONY: clean
clean: tools-clean
	rm -rf dependencies

dependencies: $(GO_MODULE_FILES)
	go mod download
	touch dependencies

.PHONY: test
test: dependencies examples/dependencies
	CGO_ENABLED=1 go test -race -count 5 -timeout 5m ./...
	(cd examples && go run -race ./simple > /dev/null)
	(cd examples && go run -race ./dag > /dev/null)

examples/dependencies: $(GO_EXAMPLES_MODULE_FILES)
	cd examples && go mod download
	touch examples/dependencies

.PHONY: lint
lint: bin/golangci-lint
	bin/golangci-lint config verify
	bin/golangci-lint run

.PHONY: tools
tools: bin/golangci-lint

.PHONY: tools-clean
tools-clean:
	rm -rf bin
	rm tools/dependencies

tools/dependencies: $(GO_TOOLS_MODULE_FILES)
	cd tools && go mod download
	touch tools/dependencies

bin/golangci-lint: tools/dependencies
	cd tools && go build -o ../bin/golangci-lint github.com/golangci/golangci-lint/v2/cmd/golangci-lint
