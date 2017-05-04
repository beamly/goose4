# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOTEST=$(GOCMD) test -v

.PHONY: default
default: updatedeps

updatedeps:
	$(GOGET)

test:
	$(GOTEST)

.PHONY: clean
clean:
	rm -rf ./goose4
