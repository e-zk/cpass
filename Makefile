# cpass (POSIX) Makefile
.POSIX:
.SUFFIXES:
.PHONY: clean build deps install

include config.mk

# by default, build
all: build

# build the binary
build: 
	go build -ldflags "-w -s" -o $(PROG) -v

# install third-party dependencies (2)
deps:
	go get golang.org/x/crypto/pbkdf2
	go get golang.org/x/crypto/ssh/terminal
	go get github.com/atotto/clipboard

# install
install:build
	install $(PROG) $(INSTALLPATH)

# clean up
clean:
	go clean
	rm -f $(PROG)
