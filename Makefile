default: unlock

unlock: ## Unlock all scripts.
	chmod -R +x ./scripts/*

install: ## Install project dependencies.
	./scripts/install

export format=image
build: ## Build application as either a Docker image or a Go binary.
	./scripts/build $(format)

export format=image
export daemon
run: ## Run application as either a Docker image, a Go binary, or a Go file.
	./scripts/run $(format) $(daemon)

export env=local
push: ## Push application's most recently built Docker image to a registry.
	./scripts/deploy $(env)

clean: ## Remove all built Go binaries.
	rm ./bin/*

ensure: ## Update dependencies.
	./scripts/install_pkgs

test: ## Run all tests.
	./scripts/run_tests