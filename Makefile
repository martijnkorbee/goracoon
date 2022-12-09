## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build: builds the command line tool to tmp/
build:
	@echo "Building GoRacoon in .tmp/ ..."
	@go build -o tmp/goracoon ./cmd/cli
	@echo "GoRacoon has been built!"

## install: runs go build and puts the binary in GOPATH bin
install:
	@echo "Installing GoRacoon in ${GOPATH}/bin ..."
	@go build -o /${GOPATH}/bin/goracoon ./cmd/cli
	@echo "GoRacoon has been installed!"
