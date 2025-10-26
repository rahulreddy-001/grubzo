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

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go
	
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

