NAME=metrics-discovery
VERSION=v0.1

all: bin

bin:
	GOARCH=amd64 GOOS=linux go build -o bin/linux/metrics-discovery
	
clean:
	rm -rf bin

test:
	go test

release: clean test bin
	hub release create -a "bin/linux/metrics-discovery" -m "$(VERSION)" "$(VERSION)"	
