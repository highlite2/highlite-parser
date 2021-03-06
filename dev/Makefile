DEV_PROJECT := $(PROJECT_NAME)-dev
DEV_COMPOSE_FILE := dev/docker-compose.yml

DEV_GOPACKAGES ?= $$(go list ./... | grep -v /vendor/ | grep -v /internal/test/mocks)
DEV_GOFILES = $$(find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: dev-rebuild
dev-rebuild:
	${INFO} "Building image..."
	docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) build import

.PHONY: dev-deps
dev-deps:
	${INFO} "Installing dependencies..."
	docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) run --rm import glide install --skip-test

.PHONY: dev-gen
dev-gen:
	${INFO} "Generating mocks..."
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) run --rm import go generate $(DEV_GOPACKAGES)

.PHONY: dev-lint
dev-lint:
	${INFO} "Running goimports..."
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) run --rm import goimports -w -local highlite2-import $(DEV_GOFILES)
	${INFO} "Running go vet..."
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) run --rm import go vet $(DEV_GOPACKAGES)
	${INFO} "Running gohint..."
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) run --rm import golint $(DEV_GOPACKAGES)

.PHONY: dev-test
dev-test:
	${INFO} "Running tests..."
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) run --rm import go test --cover -race $(DEV_GOPACKAGES)

.PHONY: dev-import
dev-import:
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) up --force-recreate -d import
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) logs -f import

.PHONY: dev-check
dev-check:
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) up --force-recreate check

.PHONY: dev-import-logs
dev-import-logs:
	@ docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) logs -f import

# Checking if config exists
ifeq ("$(wildcard dev/config.env)","")
  $(error dev/config.env file not found. Run "cp dev/_config.env dev/config.env" and override default values)
endif
