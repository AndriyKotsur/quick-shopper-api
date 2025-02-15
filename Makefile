.PHONY: build
build: ## Build the production docker image.
	docker compose build

.PHONY: start
start: ## Start the production docker container.
	docker compose up -d

.PHONY: stop
stop: ## Stop the production docker container.
	docker compose down

.PHONY: test
test: ## Run tests.
	go test ./...

.PHONY: test-unit
test-unit: ## Run unit tests.
	go test ./...

.PHONY: test-coverage
test-coverage: ## Run test coverage.
	go test -cover ./...