
test-local:
	sam build
	sam local start-api


# curl "http://localhost:3000/stori-test" -d "{\"email\":\"ente_11@hotmail.com\"}"
