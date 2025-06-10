# https://tech.davis-hansson.com/p/make/
SHELL := bash
.ONESHELL:

.PHONY: all
all: test lint build markdown

.PHONY: dummy-csv
dummy-csv:
	@echo 'AAA;my company;my city;my country' > data/file/owner.csv

.PHONY: test
test: dummy-csv
	go test ./...

.PHONY: lint
lint: dummy-csv
# See .golangci.yml
	go mod tidy -diff
	golangci-lint run

.PHONY: build
build: dummy-csv
	export CGO_ENABLED=0; go build

.PHONY: markdown
markdown: build
	rm docs/*.md
	./icm doc markdown docs

# Individual commands

.PHONY: audit
audit:
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: format
format:
	gofumpt -l -w .

.PHONY: download-owners
download-owners: build
	./icm download-owners -o data/file/owner.csv

# man-pages is also defined in goreleaser.yml
.PHONY: man-pages
man-pages: build
	mkdir -p man-pages
	./icm doc man man-pages

.PHONY: completions
completions: build
	mkdir -p completions
	./icm completion bash > completions/icm.bash
	./icm completion zsh > completions/icm.zsh
	./icm completion fish > completions/icm.fish

