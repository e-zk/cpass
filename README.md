# cpass
simple password manager written in Go.
based on the [CryptoPass Chrome extension](https://github.com/dchest/cryptopass/ "CryptoPass GitHub") and compatible with the [Android implementation](https://f-droid.org/en/packages/krasilnikov.alexey.cryptopass/ "CryptoPass Android F-Droid Page")'s JSON backup files.

the basic principal is that your password is generated from a secret, and your username/site pair:

	password = base64(pbkdf2(secret, username@url))

the password is then cut to the desired length.

note: the PBKDF2 algorithm used uses SHA-256, 5000 iterations.

after the the secret key is given, cpass will copy the password to the clipboard via the [xsel(1)](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command.

currently cpass only supports Unix-like systems (GNU/Linux, and \*BSD).

# building
using the `make(1)` command:

	$ make build

cpass uses a POSIX makefile to ensure compatibility between both GNU make and BSD make.

## dependencies

* crypto/pbkdf2:	[golang.org/x/crypto/pbkdf2](https://golang.org/x/crypto/pbkdf2)
* crypto/ssh/terminal:	[golang.org/x/crypto/ssh/terminal](golang.org/x/crypto/ssh/terminal)

cpass also depends on the [xsel(1)](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command. install it using your package manager.

# usage

	usage: cpass <command> [<args>]
	
	command can be one of:
	help			print this help message
	ls			list available bookmarks
	find <string>		find a password containing <string>
	open <user@site>	open a bookmark

in cpass, *bookmarks* are password entries, they consist of: username, site URL and length. they are listed in the following format:

	person@www.google.com (18)
	└┬───┘ └┬───────────┘  ├┘
	 │      │              └ password length
	 │      └ site URL
	 └ username

## finding bookmarks
list all available bookmarks:

	$ cpass ls
	person@www.google.com (18)
	test@site.gov (12)

filter bookmarks:

	$ cpass find site.gov
	test@site.gov (12)

	$ cpass find person
	person@www.google.com (18)

## opening bookmarks

	$ cpass open test@site.gov
	secret:
	copied to clipboard.
