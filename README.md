# cpass
simple password manager written in Go.
based on the [CryptoPass Chrome extension](https://github.com/dchest/cryptopass/ "CryptoPass GitHub") and compatible with the [Android implementation](https://f-droid.org/en/packages/krasilnikov.alexey.cryptopass/ "CryptoPass Android F-Droid Page")'s JSON backup files.

the basic principal is that your password is generated from a secret, and your username/site pair:

	password = base64(pbkdf2(secret, username@url))

the password is then cut to the desired length.

note: the PBKDF2 algorithm used uses SHA-256, 5000 iterations.

after the the secret key is given, cpass will copy the password to the clipboard via the [xsel(1)](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command.

# building
run `make`.

cpass uses a POSIX makefile to ensure compatibility between both GNU make and BSD make.

## dependencies
cpass depends on the [xsel(1)](http://www.vergenet.net/~conrad/software/xsel/ "xsel Homepage") command. install it using your package manager

# usage
`bookmarks` are password entries, they consist of: username, site URL and length. they are listed in the following format:

	person@www.google.com (18)
        |      |               L password length
	|      L site URL      
	L username

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
