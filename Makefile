run: build 
	@./bin/distributed-cache

build: 
	@go build -o bin/distributed-cache

test: 
	@go test -v ./...

runfollower: build
	@./bin/distributed-cache -listenaddr=:4000 -leaderaddr=:3000
