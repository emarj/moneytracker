# Go parameters
BINARY_NAME=./bin/moneytracker
	
all: build run
build: build-frontend build-backend
build-docker: build-frontend build-backend-linux build-docker-image
build-backend-docker: build-backend-linux build-docker-image
build-frontend:
		(cd frontend && pnpm run build)
build-backend:
		go build -o $(BINARY_NAME) -v cmd/server/main.go
build-backend-linux:
		GOOS=linux go build -o $(BINARY_NAME)_linux -v cmd/server/main.go
build-docker-image:
		docker build -t emarj/moneytracker:v2 .
#		# docker save -o bin/moneytracker_docker.tar moneytracker
test:
		go test -v ./...
clean:
		go clean
		rm -f $(BINARY_NAME)
dev:
		MT_FRONTEND_URL=https://localhost:5173/ go run ./cmd/server/main.go
push:
		docker push emarj/moneytracker
run:
		docker run -p 3245:3245 -v $(shell pwd)/data:/data emarj/moneytracker:v2
prod: 
		docker run -d -p 3245:3245 -v /home/marco/moneytracker/data:/data emarj/moneytracker:v2

