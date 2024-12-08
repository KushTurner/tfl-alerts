BINARY_NAME=tfl-alerts

build:
	go build -o bin/${BINARY_NAME} ./cmd/main.go

run:
	go run cmd/main.go

test:
	go test -v ./...

clean:
	go clean
	rm ./bin/${BINARY_NAME}

build-deploy:
	go build -o main ./cmd