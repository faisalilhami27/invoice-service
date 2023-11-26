GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Live reload:
watch-prepare: ## Install the tools required for the watch command
	go install github.com/cosmtrek/air@latest

watch: ## Run the service with hot reload
	air

## Build:
build: ## Build the service
	go build -o invoice-service

## Docker:
docker-build: ## Start the service in docker
	docker-compose up -d --build --force-recreate

## Test:
test: ## Run test
	go test -v -cover ./...

## Mock
mock-prepare: ## Install mockery
	go install github.com/vektra/mockery/v2@v2.36.0

mock: ## Generate mock files
	mockery --all --keeptree --recursive=true --outpkg=mocks

linter:
	golangci-lint run --out-format html > golangci-lint.html

migrate-file:
	migrate create -ext sql -dir migrations $(name)

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:  '
	@echo '  ${YELLOW}make${RESET} ${GREEN}<command>${RESET}'
	@echo ''
	@echo 'Commands:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)