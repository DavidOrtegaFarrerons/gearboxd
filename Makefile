.PHONY: build
build:
	docker compose up -d --build
.PHONY: test
test:
	go test ./...