# Project variables
# -----------------------------------------------------------------------------
PROJECT_NAME := highlite2-import

# Include targets
# -----------------------------------------------------------------------------
include dev/Makefile

.PHONY: check-translations
check-translations:
	@export TRANSLATIONS_FILE_PATH=data/Highlite-RUS-13-10-2018.csv && export TRANSLATIONS_FILE_SEPARATOR=";" && go run main.go -action=translation-check

.PHONY: test
test:
	go test ./...

# Common functions
# -----------------------------------------------------------------------------
# Cosmetics
YELLOW := "\e[1;33m"
NC := "\e[0m"

# Shell Functions
INFO := @bash -c '\
  printf $(YELLOW); \
  echo "=> $$1"; \
  printf $(NC)' VALUE

# Check and inspect logic
INSPECT := $$(docker-compose -p $$1 -f $$2 ps -q $$3 | xargs -I ARGS docker inspect -f "{{ .State.ExitCode }}" ARGS)

CHECK := @bash -c '\
  if [[ $(INSPECT) -ne 0 ]]; \
  then exit $(INSPECT); fi' VALUE
