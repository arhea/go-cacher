# process run during github actions
ci: install test

# install dependencies
install:
	go get .

# install and upgrade dependencies
deps: install
	go mod tidy
	go get -u ./...

# run lint using golangci-linters
lint:
	golangci-lint run

# run tests
test:
	go test -cover ./...
