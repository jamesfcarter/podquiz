fmt:
	go fmt ./...

bindata:
	go-bindata -pkg assets -o internal/assets/assets.go -prefix internal/data internal/data/views/ internal/data/static/* 

test: bindata fmt
	go test ./...

build:  bindata
	go build github.com/jamesfcarter/podquiz/cmd/podquiz
