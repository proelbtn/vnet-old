all: build

build: vnet

vnet: cmd/vnet/*
	go build -o vnet $^

.PHONY: build test vnet
