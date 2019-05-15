APP_VERSION?=latest
BUILD?=packr2 build -ldflags="-w -s"

default: build

build-all: generate format vet
	GOOS=linux GOARCH=amd64 $(BUILD) -o qmetry_uploader_linux main.go
	GOOS=darwin GOARCH=amd64 $(BUILD) -o qmetry_uploader_osx main.go
	GOOS=windows GOARCH=amd64 $(BUILD) -o qmetry_uploader_win.exe main.go
	upx qmetry_uploader_osx
	upx qmetry_uploader_linux
	upx qmetry_uploader_win.exe

build: generate format vet
	$(BUILD) -o qmetry_uploader main.go
	upx qmetry_uploader

generate:
	go generate

test:
	go test ./...

format:
	go fmt ./...

vet:
	go vet ./...
