.PHONY: worker all

build-run-worker:
	docker compose -f worker/compose.yaml up --build

down-worker:
	docker compose -f worker/compose.yaml down

