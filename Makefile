
run-db:
	docker compose up

test-local:
	sam build
	sam local start-api
