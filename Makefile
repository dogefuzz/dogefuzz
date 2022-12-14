docker.test.build:
	docker compose -f ./infrastructure/docker-compose.test.yml build

docker.test.up:
	docker compose -f ./infrastructure/docker-compose.test.yml up -d

docker.test.down:
	docker compose -f ./infrastructure/docker-compose.test.yml down -v

test.unit:
	go test ./...

test.integration:
	go test -tags=integration ./it -v -count=1
