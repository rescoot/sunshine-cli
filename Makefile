BINARY := sunshine
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X github.com/rescoot/sunshine-cli/cmd.Version=$(VERSION)"

.PHONY: build install clean test man

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

clean:
	rm -f $(BINARY)

test:
	go test ./...

man:
	go run $(LDFLAGS) . docs man --dir man

install-man: man
	install -d $(DESTDIR)$(PREFIX)/share/man/man1
	install -m 644 man/*.1 $(DESTDIR)$(PREFIX)/share/man/man1/
