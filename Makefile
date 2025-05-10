.PHONY: test
test:
	cd backend && go test -v -shuffle=on ./...

# make test-div test=○○
.PHONY: test-div
test-div:
	cd backend && go test -v ./$(test)/...

.PHONY: test-with-coverage
test-with-coverage:
	go test -v -shuffle=on ./... -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o coverage.html

.PHONY: lint
lint:
	cd backend && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && golangci-lint run