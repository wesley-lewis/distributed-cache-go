run: build 
	@./bin/distributed-cache

build: 
	@go build -o bin/distributed-cache

test: 
	@go test -v ./...
