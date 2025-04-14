
up: down
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v

swagger:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/config.yaml api/swagger.yaml


protobuf:
	protoc --go_out=./internal/grpc --go_opt=paths=source_relative \
               --go-grpc_out=./internal/grpc --go-grpc_opt=paths=source_relative \
               api/pvz.proto

test-cover:
	go test --cover ./internal/...

lint:
	golangci-lint run --config golangci.yaml