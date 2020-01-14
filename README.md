# cpass
Simple password manager written in Go, based on the [CryptoPass Chrome extension](https://github.com/dchest/cryptopass/ "CryptoPass GitHub") and compatible with the [Android implementation](https://f-droid.org/en/packages/krasilnikov.alexey.cryptopass/ "CryptoPass Android F-Droid Page")'s JSON backup files.

The basic principle is that your password is generated from a secret, and your username/site pair, then the resulting key is cut to the desired length:

	password = base64(pbkdf2(secret, username@url))[:length]

Note: The PBKDF2 algorithm used in cpass uses SHA-256, 5000 iterations.

After the the secret key is given, cpass will copy the resulting pbkdf2 key (your password) to the clipboard via the [xsel(1)](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command.

Currently cpass only supports Unix-like systems (GNU/Linux, and \*BSD).

## building
Using the `make(1)` command:

	$ # install dependencies (optional)
	$ make deps
	$ # build cpass binary
	$ make build

Note: A POSIX makefile is used to ensure compatibility between both GNU and BSD systems.

### dependencies
cpass depends on:

* crypto/pbkdf2: [golang.org/x/crypto/pbkdf2](https://golang.org/x/crypto/pbkdf2)
* crypto/ssh/terminal: [golang.org/x/crypto/ssh/terminal](golang.org/x/crypto/ssh/terminal)

You can install these dependencies by running `make deps`.

cpass also depends on the [xsel(1)](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command. You can install it using your package manager.

## usage

	usage: ./cpass [-b <path>] <command> [<args>]
	
	where:
		-b <path>		path to bookmarks file
	
	command can be one of:
		help			print this help message
		ls			list available bookmarks
		find <string>		find a password containing <string>
		open <user@site>	open bookmark

In cpass, *bookmarks* are password entries, they consist of: username, site URL and length. they are listed in the following format:

	person@www.google.com (18)
	└┬───┘ └┬───────────┘  ├┘
	 │      │              └ password length
	 │      └ site URL
	 └ username

### finding bookmarks
To list all available bookmarks, simply run `cpass ls`:

	$ cpass ls
	person@www.google.com (18)
	test@site.gov (12)

To find passwords containing a specific string; run `cpass find <string>`:

	$ cpass find site.gov
	test@site.gov (12)

cpass will also use the string to search through usernames:

	$ cpass find person
	person@www.google.com (18)

### opening bookmarks

	$ cpass open test@site.gov
	secret (will not echo):
	copied to clipboard.
