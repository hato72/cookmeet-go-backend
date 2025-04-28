.PHONY: test
test:
	cd backend && go test -v -shuffle=on ./...


.PHONY: test-with-coverage
test-with-coverage:
	go test -v -shuffle=on ./... -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o coverage.html