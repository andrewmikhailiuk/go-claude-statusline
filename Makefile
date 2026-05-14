BIN    := claude-statusline
PREFIX ?= $(HOME)/.claude

.PHONY: build install clean test fmt vet

build:
	go build -ldflags="-s -w" -o $(BIN) .

install: build
	install -m 755 $(BIN) $(PREFIX)/$(BIN)

clean:
	rm -f $(BIN)

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...
