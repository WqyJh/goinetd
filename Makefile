all: fmt build

fmt:
	go fmt ./

build:
	go build ./

clean:
	rm -rf ./goinetd
