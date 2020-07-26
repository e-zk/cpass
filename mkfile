PROG=cpass
PREFIX=/usr/local
BIN=${PREFIX}/bin

default:V: $PROG

$PROG: main.go
	go build -o $PROG -v

clean:V:
	go clean

install:V: $PROG
	install $PROG $BIN

uninstall:V:
	rm ${BIN}/${PROG}
