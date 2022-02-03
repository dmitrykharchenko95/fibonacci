build:
	go build -v -o ./bin/fibonacci ./cmd/fibonacci

run: build
	./bin/fibonacci

test:
	go test -v -race ./internal/... ./config/...

docker-build:
	sudo docker-compose build

docker-up: docker-build
	sudo docker-compose up

generate:
	mkdir -p internal/server/grpc/pb

	protoc \
    		--proto_path=proto/ \
    		--go_out=. \
    		--go-grpc_out=. \
    		proto/*.proto