_Note: cpass is a work in progress. Work is currently being done to rewrite most of it, with support for multiple password stores and encryption. As such at this time cpass should be considered unstable._

# cpass

CLI password manager written in Go.

In cpass password entries or "bookmarks" are identified in the format `user@domain` (e.g. `sam@website.org`). cpass doesn't actually store passwords anywhere. Instead, a password derivation algorithm is used (PBKDF2, SHA256, 5000 iterations) on your username+domain identifier, plus a secret key. The resulting of the key derivation algorith is encoded in base64 and stripped down to your desired length - this is your password.

*Remember: your passwords are generated using a secret key - without knowledge of this key your passwords cannot be derived.*

To build cpass, run `mk` or `make` (the `Makefile` supports both BSD and GNU make).

## usage

```console
cpass [-p] [-f] [-b FILE] COMMAND
  -p               print password to stdout.
  -b FILE          bookmarks JSON file to use.
  -f               do not ask before removing.
  COMMAND:
    help           show this help message.
    ls             list bookmarks.
    find STRING    search for bookmarks containing substring STRING.
    open BOOKMARK  open bookmark with identifier in `user@domain`
                   format.
    add ID KEYLEN  add a new bookmark. ID is `site@domain`, KEYLEN is
                   the desired limit of key length.
    rm ID          remove password. will ask to confirm unless -f is
                   given.
```

*Note: functionality for adding and removing bookmarks has not yet been implemented.*

Accessing a a bookmark is easy:

```console
$ cpass open myuser@website.com
secret (will not echo):
copied to clipboard.
```

Listing available bookmarks:

```console
$ cpass ls
joe.blogs@firefox.com (24)
jjb@tutanota.com (16)
```

Here, the number in parenthesis is the length of the password.

Finding bookmarks:

```console
$ cpass find substring
joe.blogs@substring.com (32)
```

## dependencies
cpass depends on the following Go modules:

* [github.com/atotto/clipboard](https://github.com/atotto/clipboard)
* [golang.org/x/crypto/pbkdf2](https://godoc.org/golang.org/x/crypto/pbkdf2)  
* [golang.org/x/crypto/ssh/terminal](https://godoc.org/golang.org/x/crypto/ssh/terminal)  

Additionally, on Linux and \*BSD `xsel` or `xclip` is required to copy passwords to the clipboard. You can probably install either of them using your package manager. 

## status

Currently there is only support for Linux and *BSD (X11), Windows and WSL2, because these are the operating systems I use daily.

cpass is still very much a work in progress, there is still much to be done.

## plans

* adding + editing bookmarks  
* removing bookmarks  
* encrypted `bookmarks.json` file for extra secrecy  
* configurable Terminal User Interface (TUI)  
* GUI version(?)  
