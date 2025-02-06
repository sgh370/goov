.PHONY: all test coverage coverage-html clean tidy lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=goov

all: test

test:
	$(GOTEST) -v ./...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -func=coverage.out

coverage-html: coverage
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	$(GOCLEAN)
	rm -f coverage.out coverage.html

tidy:
	$(GOMOD) tidy

lint:
	golangci-lint run

.DEFAULT_GOAL := all
