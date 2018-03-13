# Project variables
# -----------------------------------------------------------------------------
PROJECT_NAME := highlite2-import

# Include targets
# -----------------------------------------------------------------------------
include dev/Makefile

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