docker.test.build:
	docker compose -f ./infrastructure/docker-compose.test.yml build

docker.test.up:
	docker compose -f ./infrastructure/docker-compose.test.yml up -d

docker.test.down:
	docker compose -f ./infrastructure/docker-compose.test.yml down -v

test.unit:
	go test $$(go list ./... | grep -v github.com/dogefuzz/dogefuzz/test) -cover

test.unit.coverage: test.unit
	@echo "Total Coverage: $$(go tool cover -func=coverage.out | grep total: | grep -Eo '[0-9]+\.[0-9]+')%"

test.integration:
	go test -tags=integration ./it -v -count=1
