.PHONY: worker all

run-worker:
	@cd worker && go run cmd/main.go -env=../.env

run-master:
	@cd master && go run cmd/main.go -env=../.env

build-run-worker:
	@docker compose -f worker/compose.yaml up --build

down-worker:
	@docker compose -f worker/compose.yaml down

