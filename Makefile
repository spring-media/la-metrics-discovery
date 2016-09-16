PROJECT_NAME=metrics-discovery
VERSION=v0.8

all: bin

bin:
	GOARCH=amd64 GOOS=linux go build -o bin/linux/$(PROJECT_NAME)
	
clean:
	rm -rf bin

test:
	go test

release: clean test bin
	hub release create -a "bin/linux/$(PROJECT_NAME)" -m "$(VERSION)" "$(VERSION)"	
