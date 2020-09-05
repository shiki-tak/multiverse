.PHONY: build

build: go.sum
	go build -o build/simappd ./example/cmd/simappd
	go build -o build/simappcli ./example/cmd/simappcli


go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify