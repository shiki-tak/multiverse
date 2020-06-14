#!/usr/bin/make -f

build:
	go build -o build/simappd ./example/cmd/simappd
	go build -o build/simappcli ./example/cmd/simappcli

