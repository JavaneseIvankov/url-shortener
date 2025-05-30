run:
	SERVERPORT=4321 go run ./cmd/app/

lint:
	golangci-lint run
