# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=./bin/moneytracker
	
all: build run
build: build-frontend build-backend
build-backend:
		$(GOBUILD) -o $(BINARY_NAME) -v cmd/server/main.go
build-frontend:
		(cd frontend && pnpm run build)
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		./$(BINARY_NAME) -dbpath="../moneytracker_sharing.sqlite" -address="localhost"

runprod: 
		./$(BINARY_NAME) -dbpath="../../moneytracker.sqlite" -prefix="/money"

