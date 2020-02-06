# cpass
Simple password manager written in Go, based on the [CryptoPass Chrome extension](https://github.com/dchest/cryptopass/ "CryptoPass GitHub") and compatible with the [Android implementation](https://f-droid.org/en/packages/krasilnikov.alexey.cryptopass/ "CryptoPass Android F-Droid Page")'s JSON backup files.

The basic principle is that your password is generated from a username@site pair, and a  secret ("master" password):

	password = base64(pbkdf2(secret, username@url))[:length]

Note: The PBKDF2 algorithm used in `cpass` uses SHA-256 and 5000 iterations in order to be backwards compatible with both the Chrome extention and the Android application.

After the the secret key is given, cpass will copy the resulting pbkdf2 key (your password) to the clipboard via the [`xsel(1)`](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command.

Currently cpass only supports Unix-like systems (GNU/Linux, and \*BSD).

## building
Using the `make(1)` command:

	# install dependencies (optional):
	make deps
	
	# build cpass binary:
	make build

Note: A POSIX makefile is used to ensure compatibility between both GNU and BSD systems.

### dependencies
cpass depends on:

* crypto/pbkdf2: [golang.org/x/crypto/pbkdf2](https://golang.org/x/crypto/pbkdf2)
* crypto/ssh/terminal: [golang.org/x/crypto/ssh/terminal](golang.org/x/crypto/ssh/terminal)

You can install these dependencies by running `make deps`.

cpass also depends on the [`xsel(1)`](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command. You can install it using your package manager.

### installing
First open `config.mk` to confirm install location. By default cpass is installed to /usr/local/bin; you will need to run `make install` as root:

	# first obtain root shell via doas(1) or sudo(1)...
	# install cpass to /usr/local:
	make install

## usage

	usage: ./cpass [-b <path>] <command> [<args>]
	
	where:
		-b <path>		path to bookmarks file
	
	command can be one of:
		help			print this help message
		ls			list available bookmarks
		find <string>		find a password containing <string>
		open <user@site>	open bookmark

In cpass, _bookmarks_ are account entries. A bookmark consists of: a username, a site URL, and the length of the password. Bookmarks are listed in the following format:

	person@www.google.com (18)
	└┬───┘ └┬───────────┘  ├┘
	 │      │              └ password length
	 │      └ site URL
	 └ username

A "bookmarks file" is a simple JSON file that holds a collection of bookmarks. An example bookmarks file, containing two bookmarks, would look like this:

	[
		{
			"url": "www.google.com",
			"username": "person",
			"length": 18
		},
		{
			"url": "site.gov",
			"username": "test",
			"length": 12
		}
	]

A user's primary bookmarks file is located at `$HOME/.CryptopassBookmarks.txt`. If a different bookmarks file is to be parsed, the path to it can be specified with the `-b` flag.

### finding bookmarks
To list all available bookmarks, simply run `cpass ls`:

	$ cpass ls
	person@www.google.com (18)
	test@site.gov (12)

To find passwords containing a specific string; run `cpass find <string>`:

	$ cpass find site.gov
	test@site.gov (12)

Note: `cpass find <string>` tries to match the given string within the whole bookmark identifier ('username@site'). So, using an entire or partial bookmark identifier works:

	$ cpass find son@wwww.google
	person@www.google.com (18)

If you wish, the output of `cpass ls` can be piped into other programs, such as `grep(1)`:

	$ cpass ls | grep -E '.+\.gov'
	test@site.gov (12)'

### opening bookmarks
To open a bookmark, simply supply cpass with your account in the 'username@site' format:

	$ cpass open test@site.gov
	secret (will not echo): 
	copied to clipboard.
