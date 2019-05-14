APP_VERSION?=latest
BUILD?=packr2 build -ldflags="-w -s"

default: build

build-all: format vet
	GOOS=windows GOARCH=amd64 $(BUILD) -o qmetry_uploader_win.exe main.go
	GOOS=linux GOARCH=amd64 $(BUILD) -o qmetry_uploader_linux main.go
	GOOS=darwin GOARCH=amd64 $(BUILD) -o qmetry_uploader_osx main.go
	upx qmetry_uploader_osx
	upx qmetry_uploader_win.exe
	upx qmetry_uploader_linux

build: format vet
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
