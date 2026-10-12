# meth - build and install
#
# go-sqlite3 is a cgo package. With CGO_ENABLED=0 it still compiles and then
# fails at runtime with "Binary was compiled with 'CGO_ENABLED=0' ... This is
# a stub", so pin it here rather than inherit a broken value from the shell.
export CGO_ENABLED = 1

BINARY := meth
PKG    := ./cmd/meth

# install target, override with: make install PREFIX=/usr/local
PREFIX ?= $(HOME)/.local
BINDIR := $(PREFIX)/bin

.PHONY: all build install uninstall test vet fmt clean

all: build

build:
	go build -o bin/$(BINARY) $(PKG)

install: build
	install -d $(BINDIR)
	install -m 0755 bin/$(BINARY) $(BINDIR)/$(BINARY)
	@echo "installed $(BINDIR)/$(BINARY)"
	@case ":$$PATH:" in \
		*":$(BINDIR):"*) ;; \
		*) echo; \
		   echo "warning: $(BINDIR) is not on your PATH, so '$(BINARY)' will not be found."; \
		   echo "add this to ~/.zshrc:"; \
		   echo "    export PATH=\"$(BINDIR):\$$PATH\"";; \
	esac

uninstall:
	rm -f $(BINDIR)/$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l -w .

clean:
	rm -rf bin
