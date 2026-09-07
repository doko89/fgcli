.PHONY: build vet test
build:
	go build -o bin/fgcli ./cmd/fgcli
vet:
	go vet ./...
test:
	go test ./...
