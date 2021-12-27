## works fine on Linux, Build will work on Windows too, the other commands have not be adapted to Windows yet.
BINARY_NAME=gosnake

.ONESHELL:
build:
ifeq ($(OS),Windows_NT)
	@SET GOARCH=amd64
	@SET GOOS=darwin
	@go build -o ./bin/${BINARY_NAME}.app ./cmd/gosnake
	@SET GOARCH=amd64
	@SET GOOS=linux
	@go build -o ./bin/${BINARY_NAME} ./cmd/gosnake
	@SET GOARCH=amd64
	@SET GOOS=windows
	@go build -o ./bin/${BINARY_NAME}.exe ./cmd/gosnake
else
	@GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}.app ./cmd/gosnake
	@GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME} ./cmd/gosnake
	@GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}.exe ./cmd/gosnake
endif


run:
	go run ./cmd/gosnake

clean:
	@go clean
	@rm ./bin/${BINARY_NAME}.app
	@rm ./bin/${BINARY_NAME}
	@rm ./bin/${BINARY_NAME}.exe

install-linters: ## Install linters ## Copied from github.com/skycoin/skycoin
	# Turn off go module when install the vendoercheck, otherwise the installation
	# will pollute the go.mod file.
	GO111MODULE=off go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	# go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	# Change to use go get -u with version when go is v1.12+
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.43.0
	# added to install goimports for make format:
	go install golang.org/x/tools/cmd/goimports@latest

lint: ## Run linters. Use make install-linters first. ## Copied from github.com/skycoin/skycoin
	GO111MODULE=off vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	@# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

mock:
	mockery --all --dir cmd --dir pkg

test:
	go test `go list ./... | grep -v mocks` -cover

format: ## Formats the code. Must have goimports installed (use make install-linters). ## Copied from github.com/skycoin/skycoin
	goimports -w -local github.com/skycoin/skycoin ./cmd
	goimports -w -local github.com/skycoin/skycoin ./pkg