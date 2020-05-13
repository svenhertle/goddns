TARGET=goddns

all: build-x64-linux

prepare:
	mkdir -p ./build

clean:
	rm -rf ./build

build-x64-linux: prepare
	GOOS=linux GOARCH=amd64 go build -o build/$(TARGET)-x64-linux .