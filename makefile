.POSIX:
.SUFFIXES:
.PHONY: clean

PREFIX = /usr/local

cpass: main.go
	go build -ldflags "-w -s" -o cpass -v main.go

install: cpass
	install -c -s -m 0755 cpass $(PREFIX)/bin

clean:
	rm -f cpass
	go clean
	rm -f cpass
