# cpass
Simple password manager written in Go, based on the [CryptoPass Chrome extension](https://github.com/dchest/cryptopass/ "CryptoPass GitHub") and compatible with the [Android implementation](https://f-droid.org/en/packages/krasilnikov.alexey.cryptopass/ "CryptoPass Android F-Droid Page") JSON backup files.

The basic principle is that your password is generated from a username@site pair, and a  secret ("master" password):

	password = base64(pbkdf2(secret, username@url))[:length]

Note: The PBKDF2 algorithm used in `cpass` uses SHA-256 and 5000 iterations in order to be backwards compatible with both the Chrome extention and the Android application.

After the the secret key is given `cpass` will copy the generated key (your password) to the clipboard via the [`xsel(1)`](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command on \*nix, or via `clip.exe` on windows.

`cpass` works on Unix-like systems (Linux, *BSD) and also Windows and WSL.

## usage

	cpass [-p] [-b path] [help|ls|find|open] [args]

### finding bookmarks
In cpass, _bookmarks_ are account entries. A bookmark consists of: a username, a site URL, and the length of the password.  
To list all available bookmarks, simply run `cpass ls`:

	$ cpass ls
	person@www.google.com (18)
	test@site.gov (12)

Bookmarks are listed in the following format:

	person@www.google.com (18)
	└┬───┘ └┬───────────┘  ├┘
	 │      │              └ password length
	 │      └ site URL
	 └ username

To find passwords containing a specific string; run `cpass find <string>`:

	$ cpass find site.gov
	test@site.gov (12)

The `find` command tries to match the given string within the whole bookmark identifier ('username@site'). So, using an entire or partial bookmark identifier works:

	$ cpass find son@wwww.google
	person@www.google.com (18)

If you wish, the output of `cpass ls` can be piped into other programs, such as `grep(1)`:

	$ cpass ls | grep -E '.+\.gov'
	test@site.gov (12)'

### opening bookmarks
To open a bookmark supply cpass with your account in the 'username@site' format:

	$ cpass open test@site.gov
	secret (will not echo): 
	copied to clipboard.

### file format
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

A user's primary bookmarks file is located at `$HOME/.CryptopassBookmarks.txt`. To use a different bookmarks file give it's path using the `-b` flag.

## building
Using the `make(1)` command:

	# install dependencies (optional)
	make deps
	
	# build cpass binary
	make build

### dependencies
cpass depends on:

* crypto/pbkdf2: [golang.org/x/crypto/pbkdf2](https://golang.org/x/crypto/pbkdf2)
* crypto/ssh/terminal: [golang.org/x/crypto/ssh/terminal](golang.org/x/crypto/ssh/terminal)

You can install these dependencies by running `make deps`.

cpass also depends on the [`xsel(1)`](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command. You can install it using your package manager.

### installing
First open `config.mk` to confirm install location. By default cpass is installed to /usr/local/bin; you will need to run `make install` as root:

	# first obtain root shell via doas(1) or sudo(1)...
	# install cpass to /usr/local
	make install

Installation location can also be changed through make flags:

	# install to $HOME/bin
	make PREFIX=$HOME install

