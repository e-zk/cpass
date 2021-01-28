.POSIX:
.SUFFIXES:
.PHONY: clean build

build:
	go build -ldflags "-w -s" -o cpass -v ./cmd/cpass

clean:
	go clean
	rm -f cpass
