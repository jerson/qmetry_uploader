APP_VERSION?=latest
PACKAGER?=packr2
BUILD?=packr2 build -ldflags="-w -s"
NAME?=qmetry_uploader
UPX?=upx

default: build

build-all: generate format vet
	GOOS=linux GOARCH=amd64 $(BUILD) -o $(NAME)_linux main.go
	GOOS=darwin GOARCH=amd64 $(BUILD) -o $(NAME)_osx main.go
	GOOS=windows GOARCH=amd64 $(BUILD) -o $(NAME)_win.exe main.go
	$(UPX) $(NAME)_osx
	$(UPX) $(NAME)_linux
	$(UPX) $(NAME)_win.exe

build-win: generate format vet
	$(BUILD) -o $(NAME).exe main.go
	$(UPX) $(NAME).exe

build: generate format vet
	$(BUILD) -o $(NAME) main.go
	$(UPX) $(NAME)

clean:
	$(PACKAGER) clean
	rm -rf assets/*.zip
	rm -rf $(NAME)
	rm -rf $(NAME)*

generate:
	go generate

test:
	go test ./...

format:
	go fmt ./...

vet:
	go vet ./...
