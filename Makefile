BINARY_NAME=tfl-alerts

build:
	go build -o bin/${BINARY_NAME} main.go

run:
	go run main.go

test:
	go test -v ./...

clean:
	go clean
	rm ./bin/${BINARY_NAME}

build-deploy:
	go build -o main .