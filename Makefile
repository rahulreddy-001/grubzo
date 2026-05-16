GOLANGCI_LINT ?= golangci-lint
GOLANGCI_LINT_INSTALL_VERSION ?= v2.11.3
PRETTIER ?= ./frontend/node_modules/.bin/prettier
GO_FILES := $(shell find ./cmd ./internal -type f -name '*.go' -print)
GO_LINT_CACHE ?= /tmp/go-build
GOLANGCI_LINT_CACHE ?= /tmp/golangci-lint-cache

.PHONY: docker-up docker-down run-react run migrate build watch lint lint-go lint-go-fmt lint-go-style lint-js lint-js-prettier install-hooks install-lint

docker-up:
	@sudo docker compose up -d

docker-down:
	@sudo docker compose down

run-react:
	@echo "Starting frontend..."
	@cd frontend && npm run dev > ../tmp/frontend.log 2>&1 &
	@echo "Frontend started (logs in tmp/frontend.log)"

run: run-react
	@go run cmd/api/main.go serve

migrate:
	@go run cmd/api/main.go migrate

build:
	@echo "Building..."
	@go build -o grubzo cmd/api/main.go

lint: lint-go lint-go-style

lint-go:
# 	@GOCACHE=$(GO_LINT_CACHE) go test ./...

lint-go-style: lint-go-fmt
# 	@GOCACHE=$(GO_LINT_CACHE) go test -run='^$$' ./...

lint-go-fmt:
	@UNFORMATTED=$$(gofmt -l $(GO_FILES)); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "Formatting Go files with gofmt..."; \
		gofmt -w $$UNFORMATTED; \
	fi

lint-js: lint-js-prettier

lint-js-prettier:
	@if [ ! -x "$(PRETTIER)" ]; then \
		echo "Prettier is not installed in frontend."; \
		echo "Run: cd frontend && npm install -D prettier"; \
		exit 1; \
	fi
	@$(PRETTIER) --check "frontend/**/*.{js,jsx,ts,tsx}" "frontend/*.{js,json,md,yaml,yml}"

install-hooks:
	@git config core.hooksPath .githooks
	@echo "Git hooks path set to .githooks"

install-lint:
	@command -v curl >/dev/null 2>&1 || { \
		echo "curl is required to install golangci-lint"; \
		exit 1; \
	}
	@curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_INSTALL_VERSION)
	@$(GOLANGCI_LINT) version

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

containers:
	docker build --target frontend_builder -t grubzo-frontend .
	docker build --target backend_builder -t grubzo-backend .
