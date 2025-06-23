.PHONY: model model-sql api build docker-build

model:
	goctl model mysql datasource --url="root:root@tcp(localhost:3306)/volctrain" --table="user" --dir ./api/models

model-sql:
	goctl model mysql ddl --src model/mysql/volctrain.sql --dir ./api/models

api:
	goctl api go -api backend.api -dir ./

build:
	GOOS=linux GOARCH=amd64 go build -o bin/backend main.go

docker-build:
	docker build --platform linux/amd64 --push -t volctrain-backend:v1.0.0 .