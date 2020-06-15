.PHONY: build

build:
	go build -mod readonly -o build/simappd ./example/cmd/simappd
	go build -mod readonly -o build/simappcli ./example/cmd/simappcli

initialize: build
	./scripts/initialize.sh
