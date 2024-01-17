.PHONY: build run test docker-build docker-run

build:
	go build -o main ./cmd/proto/main.go

run: build
	./main

test:
	go test -v -count=1 ./...

docker-build:
	docker build -t myapp .

docker-run: docker-build
	docker run -p 50051:50051 myapp
