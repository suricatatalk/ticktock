# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
	
all: build
build: 
	$(GOBUILD) -v ./...
clean: 
	$(GOCLEAN)
test: 
	$(GOTEST) -v ./...
