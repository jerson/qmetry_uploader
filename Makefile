APP_VERSION?=latest
BUILD?=go build -ldflags="-w -s"

default: build

build: format vet
	$(BUILD) -o qmetry_uploader main.go
	upx -9 qmetry_uploader

test:
	go test $$(go list ./... | grep -v /vendor/)

format:
	go fmt $$(go list ./... | grep -v /vendor/)

vet:
	go vet $$(go list ./... | grep -v /vendor/)
