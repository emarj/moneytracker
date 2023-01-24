# Go parameters
BINARY_NAME=./bin/moneytracker
	
all: build run
build: build-frontend build-backend
build-docker: build-frontend build-backend-linux build-docker-image
build-backend-docker: build-backend-linux build-docker-image
build-frontend:
		(cd frontend && pnpm run build)
build-backend:
		./build.sh build $(BINARY_NAME)
build-backend-linux:
		GOOS=linux GOARCH=amd64 ./build.sh build $(BINARY_NAME)
build-docker-image:
		./build_docker.sh
build-docker-image-compile:
		docker build -t emarj/moneytracker:v2 -f Dockerfile.compile .
test:
		go test -v ./...
clean:
		go clean
		rm -f $(BINARY_NAME)
dev:
		MT_FRONTEND_URL=http://localhost:5173/ go run ./cmd/server/main.go --local
dev-temp:
		MT_FRONTEND_URL=http://localhost:5173/ go run ./cmd/server/main.go --local --populate --tempDB
dev-no-proxy:
		go run ./cmd/server/main.go
push:
		docker push -a emarj/moneytracker
docker-run:
		docker run --rm -p 3245:3245 -v $(shell pwd)/data:/data emarj/moneytracker:latest
docker-prod: 
		docker run -d -p 3245:3245 -v /home/marco/data:/data emarj/moneytracker:latest
gen:
		jet -source=sqlite -dsn="./data/moneytracker.sqlite" -path="./.gen"
		rm -rf .gen/model
schema:
		rm data/moneytracker.sqlite* || true
		sqlite3 data/moneytracker.sqlite < store/sqlite/schema.sql
populate:
		go run ./cmd/populate

