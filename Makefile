
.PHONY: deps
deps:
	${INFO} "Installing dependencies"
	@ glide install

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