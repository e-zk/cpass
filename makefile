.POSIX:
.SUFFIXES:
.PHONY: clean build

build:
	go build -ldflags "-w -s" -o cpass -v main.go

clean:
	go clean
