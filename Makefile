GOCMD=go
BINARY_NAME=dtree

build:
	$(GOCMD) build -v -o $(BINARY_NAME) ./cmd/...

install:
	$(GOCMD) install ./cmd/...

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)
