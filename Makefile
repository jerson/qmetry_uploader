APP_VERSION?=latest
BUILD?=packr2 build -ldflags="-w -s"
NAME?=qmetry_uploader

default: build

build-all: generate format vet
	GOOS=linux GOARCH=amd64 $(BUILD) -o $(NAME)_linux main.go
	GOOS=darwin GOARCH=amd64 $(BUILD) -o $(NAME)_osx main.go
	GOOS=windows GOARCH=amd64 $(BUILD) -o $(NAME)_win.exe main.go
	upx $(NAME)_osx
	upx $(NAME)_linux
	upx $(NAME)_win.exe

build-win: generate format vet
	$(BUILD) -o $(NAME).exe main.go
	upx $(NAME).exe

build: generate format vet
	$(BUILD) -o $(NAME) main.go
	upx $(NAME)

generate:
	go generate

test:
	go test ./...

format:
	go fmt ./...

vet:
	go vet ./...
