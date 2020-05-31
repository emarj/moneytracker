# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=./Bin/moneytracker
	
all: build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		 DBPATH="./db/moneytracker.sqlite" ./$(BINARY_NAME)

prod: 
		DBPATH="../db/moneytracker.sqlite" PREFIX="/money" ./$(BINARY_NAME)
