# binary name
BIN = cpass

# install location
PREFIX = /usr/local
INSTALLPATH = $(PREFIX)/bin/$(BIN)

# go
GOCMD = go
GOBUILD = $(GOCMD) build
GOGET = $(GOCMD) get
