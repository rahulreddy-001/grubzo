GOLANGCI_LINT ?= go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3
GOLANGCI_LINT_INSTALL_VERSION ?= v2.11.3
PRETTIER ?= ./frontend/node_modules/.bin/prettier
GOCACHE ?= /tmp/go-build
GOLANGCI_LINT_CACHE ?= /tmp/golangci-lint-cache
GO_PACKAGES := ./cmd/... ./internal/...
GO_FILES := $(shell find ./cmd ./internal -type f -name '*.go' -print)

.PHONY: docker-up docker-down run-react run migrate build watch
		install-hooks lint test

docker-up:
	@sudo docker compose up -d

docker-down:
	@sudo docker compose down

run:
	@go run cmd/api/main.go serve

migrate:
	@go run cmd/api/main.go migrate

build:
	@echo "Building..."
	@go build -o grubzo cmd/api/main.go

watch:run-react docker-up
	@if command -v air > /dev/null; then \
		air; \
	else \
		read -p "Go's 'air' is not installed. Install? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi \
	fi

install-hooks:
	@git config core.hooksPath .githooks
	@echo "Git hooks path set to .githooks"

lint:
	@if [ "$$(gofmt -l $(GO_FILES) | wc -l)" -gt 0 ]; then \
		echo "The following files are not formatted:"; \
		gofmt -l $(GO_FILES); \
		exit 1; \
	fi
	@GOCACHE=$(GOCACHE) go vet $(GO_PACKAGES)
	@GOCACHE=$(GOCACHE) GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) $(GOLANGCI_LINT) run --allow-serial-runners $(GO_PACKAGES)

test:
	@GOCACHE=$(GOCACHE) go test -race -covermode=atomic -coverprofile=coverage.out $(GO_PACKAGES)

REPO ?= rohan001
IMAGE_NAME ?= grubzo-backend
TAG ?= $(shell git rev-parse --short HEAD)

docker-build:
	docker build -t $(REPO)/$(IMAGE_NAME):$(TAG) .

docker-push:
	docker push $(REPO)/$(IMAGE_NAME):$(TAG)
