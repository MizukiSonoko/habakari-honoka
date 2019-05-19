
build:
	go build -o honoka

bench:
	go test -test.bench=".*" ./...

honoka: build
	CGO_ENABLED=1 CC=gcc go test -test.bench=".*" ./... | ./honoka

