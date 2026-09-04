GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test ./...

.PHONY: build test clean

build:
	$(GOBUILD) ./...

test:
	$(GOTEST)

clean:
	$(GOCMD) clean
