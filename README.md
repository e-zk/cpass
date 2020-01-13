# cpass
simple password manager written in Go.
based on the [CryptoPass Chrome extension](https://github.com/dchest/cryptopass/ "CryptoPass GitHub") and compatible with the [Android implementation](https://f-droid.org/en/packages/krasilnikov.alexey.cryptopass/ "CryptoPass Android F-Droid Page")'s backups.

the basic principal is that your password is generated from a secret, and your username/site pair:

	password = base64(pbkdf2(secret, username@url))

the password is then cut to the desired length.

note: the PBKDF2 algorithm used uses SHA-256, 5000 iterations.

# building
run `make`.

cpass uses a POSIX makefile to ensure compatibility between both GNU make and BSD make.

# usage

