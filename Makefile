GOCMD=go
BINARY_NAME=dtree

build:
	$(GOCMD) build -v -o $(BINARY_NAME) ./cmd/...

install:
	$(GOCMD) install ./cmd/...

test:
	$(GOCMD) test -v -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)
	rm -f cover*
