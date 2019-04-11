
build:
	go build -o honoka

bench:
	go test -test.bench=".*" ./...

honoka: build
	go test -test.bench=".*" ./... | ./honoka

