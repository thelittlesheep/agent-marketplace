BINARY := status-line
INSTALL_DIR := $(HOME)/.local/bin

.PHONY: build test install clean

build:
	go build -o $(BINARY) .

test:
	go test -v ./...

install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed $(BINARY) to $(INSTALL_DIR)"

clean:
	rm -f $(BINARY)
