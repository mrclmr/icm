BIN_DIR := $(GOPATH)/bin
BUILD_DIR := build
# man-pages is also defined in goreleaser.yml
MAN_DIR := man-pages
DOCS_DIR := docs
BINARY := icm

.PHONY: all
all: test lint build man markdown

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint: test
	golangci-lint run --enable golint --enable goimports

.PHONY: build
build: lint
	export CGO_ENABLED=0; go build -o $(BUILD_DIR)/$(BINARY)

.PHONY: man
man: build
	$(shell $(BUILD_DIR)/$(BINARY) misc man $(MAN_DIR)/man1)

.PHONY: markdown
markdown: build
	$(shell $(BUILD_DIR)/$(BINARY) misc markdown $(DOCS_DIR))

.PHONY: install
install: man markdown
	cp $(BUILD_DIR)/$(BINARY) $(BIN_DIR)/$(BINARY)

.PHONY: clean
clean:
	go clean -x -testcache
	rm -rf $(BUILD_DIR)

.PHONY: fmt
fmt:
	goimports -w $(shell find . -type f -name '*.go')