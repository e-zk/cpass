# cpass (POSIX) Makefile
.POSIX:
.SUFFIXES:

.include <config.mk>

# by default, build
all: build

# build the binary
build: 
	$(GOBUILD) -o $(BIN) -v 

# install third-party dependencies (2)
deps:
	$(GOGET) golang.org/x/crypto/pbkdf2
	$(GOGET) golang.org/x/crypto/ssh/terminal

# install
install:build
	[ -f $(INSTALLPATH) ] && rm -i $(INSTALLPATH)
	cp -v $(BIN) $(INSTALLPATH)
	chmod 2555 $(INSTALLPATH)

# clean up
clean:
	$(GOCMD) clean
	rm -f $(BIN)
