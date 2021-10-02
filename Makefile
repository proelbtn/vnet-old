all: build

build: vnet

vnet: cmd/vnet/main.go
	go build -o vnet $^

install:
	sudo cp ./vnet /usr/local/bin

.PHONY: build test vnet
