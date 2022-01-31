build:
	go build -v -o ./bin/fibonacci ./cmd/fibonacci

run: build
	redis-server --port 6379
	./bin/fibonacci

test:
	go test -v -race ./internal/...

generate:
	mkdir -p internal/server/grpc/pb

	protoc \
    		--proto_path=proto/ \
    		--go_out=. \
    		--go-grpc_out=. \
    		proto/*.proto