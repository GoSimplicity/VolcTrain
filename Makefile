docker-build:
	docker build -t Bamboo/gomodd:v1.22.1 .

docker-start:
	docker-compose -f docker-compose.yaml up -d

docker-stop:
	docker-compose -f docker-compose.yaml down

docker-net-remove:
	docker network rm volctrain_net

dev: docker-build docker-start

stop: docker-stop docker-net-remove