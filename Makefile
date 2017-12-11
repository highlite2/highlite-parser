GOPACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: install-tools
install-tools:
	${INFO} "Installing tools for development..."
	@ go get -u github.com/golang/lint/golint
	@ go get -u golang.org/x/tools/cmd/goimports

.PHONY: deps
deps:
	${INFO} "Installing dependencies..."
	@ glide install

.PHONY: test
test:
	${INFO} "Running tests..."
	@ go test -v ${GOPACKAGES}

.PHONY: check
check:
	${INFO} "Running goimports..."
	@ goimports -w -local highlite-parser $(GOFILES)

	${INFO} "Running go vet..."
	@ go vet ${GOPACKAGES}

	${INFO} "Running gohint..."
	@ golint ${GOPACKAGES}

	${INFO} "DONE"


# ======================================================================================================================
# COMMON FUNCTIONS
# ======================================================================================================================

# Cosmetics
YELLOW := "\e[1;33m"
NC := "\e[0m"

# Shell Functions
INFO := @bash -c '\
  printf $(YELLOW); \
  echo "=> $$1"; \
  printf $(NC)' VALUE