gen-swagger:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/config.yaml api/swagger.yaml

up: down
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v
