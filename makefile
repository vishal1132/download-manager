.DEFAULT_GOAL=help
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make [target...]"
	@echo ""
	@echo "Useful commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
run: ## to run the downloadmanager
	@go run main.go

build: ## to build the binary for filemanager
	@go build  -o vdl .

deploy: ## to deploy the binary to bin directory for unix and linux like os
	#move the output file to /usr/local/bin

buildanddeploy: build deploy

	