.PHONY: build

build: go.sum
	go build -o build/simd ./example/cmd/simd
	go build -o build/simcli ./example/cmd/simcli


go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify