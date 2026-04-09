.PHONY: build
build:
	docker compose up -d --build
.PHONY: seed
seed:
	docker compose exec app ./seed
.PHONY: test
test:
	go test ./...