GO=/usr/local/go/bin/go
BINARY=blog

.PHONY: build run clean tidy

build:
	$(GO) build -o $(BINARY) main.go

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

tidy:
	$(GO) mod tidy

image:
	docker build -t registry.gitlab.com/nmcostello/blog .

push:
	docker push registry.gitlab.com/nmcostello/blog
