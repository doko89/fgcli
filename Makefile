.PHONY: build vet test install
build:
	go build -o bin/fgcli ./cmd/fgcli
vet:
	go vet ./...
test:
	go test ./...
# Install to ~/.local/bin (add it to PATH once: export PATH="$$HOME/.local/bin:$$PATH")
install: build
	mkdir -p $(HOME)/.local/bin
	install -m755 bin/fgcli $(HOME)/.local/bin/fgcli
