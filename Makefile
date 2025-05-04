.PHONY: build tidy

build:
	docker run -it -v "$(PWD):/app" -w /app krakend/builder:2.9.4 go build -buildmode=plugin -o krakend-grpc-proxy.so .

tidy:
	docker run -it -v "$(PWD):/app" -w /app krakend/builder:2.9.4 go mod tidy

test:
	go test ./... -v