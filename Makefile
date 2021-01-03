GOCMD=go
BINARY_NAME=dtree

build:
	$(GOCMD) build -v -o $(BINARY_NAME)

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)
