APP_VERSION?=latest
BUILD?=go build -ldflags="-w -s"

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

test:
	go test $$(go list ./... | grep -v /vendor/)

format:
	go fmt $$(go list ./... | grep -v /vendor/)

vet:
	go vet $$(go list ./... | grep -v /vendor/)
