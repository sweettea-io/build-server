default: unlock

unlock: ## Unlock all scripts
	chmod -R +x ./scripts/*

install: ## Install project dependencies
	./scripts/install