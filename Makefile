# cpass (POSIX) Makefile
.POSIX:
.SUFFIXES:

# macros
BIN = cpass
GOCMD = go
GOBUILD = $(GOCMD) build
GOGET = $(GOCMD) get

# by default, build
all: build

# build the binary
build: 
	$(GOBUILD) -o $(BIN) -v 

# install third-party dependencies (2)
deps:
	$(GOGET) golang.org/x/crypto/pbkdf2
	$(GOGET) golang.org/x/crypto/ssh/terminal

# clean up
clean:
	$(GOCMD) clean
	rm -f $(BIN)
