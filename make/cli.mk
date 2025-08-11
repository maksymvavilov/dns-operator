.PHONY: build-cli
build-cli: CLI_VERSION ?= dev
build-cli: fmt vet
	@set -e; \
	GIT_SHA=$$(git rev-parse --short=7 HEAD 2>/dev/null) || { \
		GIT_HASH=$${GITHUB_SHA:-NO_SHA}; \
	}; \
	if [ -z "$$GIT_HASH" ]; then \
		GIT_DIRTY=$$(git diff --stat); \
		if [ -n "$$GIT_DIRTY" ]; then \
			GIT_HASH=$${GIT_SHA}-dirty; \
		else \
			GIT_HASH=$${GIT_SHA}; \
		fi; \
	fi; \
	LDFLAGS="-X 'main.gitSHA=$$GIT_HASH' -X 'main.version=$(CLI_VERSION)'"; \
	go build -ldflags "$$LDFLAGS" -o bin/kubectl-dns cmd/plugin/*.go
	@echo "To embed plugin in kubectl add ./bin to your PATH"