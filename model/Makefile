generate-mocks: ## Generate Mocks
generate-mocks:
	mockery --dir device --name Device
	mockery --dir library --name Library
	mockery --dir library --name Media
	mockery --dir manager --name Manager
	mockery --dir session --name Session

test: ## Run tests
test:
	go test -v -count=1 -race ./...
