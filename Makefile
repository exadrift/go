GO_VERSION ?= 1.26
LINTER_VERSION ?= 2.12.2
GO_IMAGE := golang:$(GO_VERSION)-alpine
GO_RUN := docker run --rm -e HOME=$$HOME -v $$HOME:$$HOME -u $(shell id -u):$(shell id -g) -v $(shell pwd):/build -w /build $(GO_IMAGE) go

.PHONY: test
test:
	$(GO_RUN) test -cover -p 1 --timeout 10m ./...

.PHONY: lint-check
lint-check:
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v$(LINTER_VERSION) golangci-lint run
