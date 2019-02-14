# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=./Bin/moneytracker
	
all: build run
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		./$(BINARY_NAME) -dbpath="../moneytracker_sharing.sqlite" -address="localhost"

runprod: 
		./$(BINARY_NAME) -dbpath="../../moneytracker.sqlite" -prefix="/money"

