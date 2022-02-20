fmt:
	go fmt ./...

test: fmt
	go test ./...

build:
	go build -tags netgo github.com/jamesfcarter/podquiz/cmd/podquiz

run: build
	./podquiz
