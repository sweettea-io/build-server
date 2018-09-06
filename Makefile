default: unlock

unlock: ## Unlock all scripts.
	chmod -R +x ./scripts/*

install: ## Install project dependencies.
	./scripts/install

export env=local
export format=image
export version=latest
build: ## Build the application for a specific environment tier as either a Docker image or a Go binary.
	./scripts/build $(env) $(format) $(version)

export env=local
export format=image
run: ## Run the application as either a Docker image (for a specific environment), a Go binary, or a Go file.
	./scripts/run $(env) $(format)

export env=local
export version=latest
push: ## Push application's most recently built Docker image (for a specific environment) to a registry.
	./scripts/push $(env) $(version)

clean: ## Remove all built Go binaries.
	rm ./bin/*

ensure: ## Update dependencies.
	./scripts/install_pkgs

test: ## Run all tests.
	./scripts/run_tests