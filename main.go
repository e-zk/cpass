package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	iterations   = 5000                       // pbkdf2 iterations
	xselPath     = "xsel"                     // path to xsel(1)
	secretPrompt = "secret (will not echo): " // secret prompt
)

var (
	xselArgs = []string{xselPath, "-i"} // xel(1) arguments
)

// A bookmark is defined as follows in the JSON backup format...
type Bookmark struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Length   int    `json:"length"`
}

// Bookmarks is a list of type Bookmark
type Bookmarks []Bookmark

// Prints program usage information
func usage() {

	fmt.Printf("usage: %s [-b <path>] <command> [<args>]\n\n", os.Args[0])
	fmt.Printf("where:\n")
	fmt.Printf("\t-b <path>\t\tpath to bookmarks file\n")
	fmt.Printf("\n")
	fmt.Printf("command can be one of:\n")
	fmt.Printf("\thelp\t\t\tprint this help message\n")
	fmt.Printf("\tls\t\t\tlist available bookmarks\n")
	fmt.Printf("\tfind <string>\t\tfind a password containing <string>\n")
	fmt.Printf("\topen <user@site>\topen bookmark\n")
}

// Generate password from given secret, and Bookmark
func genPassword(secret []byte, bmark *Bookmark) string {

	// the salt is "user@site"
	salt := fmt.Sprintf("%s@%s", bmark.Username, bmark.Url)

	// generate the pbkdf2 key based on the input values
	pbkdf2Hmac := pbkdf2.Key([]byte(secret), []byte(salt), iterations, 32, sha256.New)

	// encode the resulting pbkdf2 key in base64
	b64Encoded := base64.StdEncoding.EncodeToString([]byte(pbkdf2Hmac))

	// cut the encoded key down to the given length; this is the final password
	return b64Encoded[:bmark.Length]
}

// Get secret from user
func inputSecret() ([]byte, error) {

	// prompt for input
	fmt.Printf(secretPrompt)

	secret, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")
	return secret, err
}

// Filter a list of bookmarks, based on a keyword and return bookmarks which match this keyword
func filterList(bmarks Bookmarks, filter string) (outBmarks Bookmarks) {

	// if the filter string is within the Username or URL
	for _, bmark := range bmarks {
		fullname := fmt.Sprintf("%s@%s", bmark.Username, bmark.Url)
		if strings.Contains(fullname, filter) {
			// append the bookmark to outBmarks
			outBmarks = append(outBmarks, bmark)
		}
	}

	// return the filtered list of bookmarks
	return outBmarks
}

// Print a list of bookmarks out
func list(bmarks Bookmarks) {

	// foreach bookmark in the bookmarks array...
	for _, bmark := range bmarks {
		fmt.Printf("%s@%s (%d)\n", bmark.Username, bmark.Url, bmark.Length)
	}
}

// Copies a string to the clipboard via xsel(1)
func clipboard(input string) {

	// xsel command
	clipCmd := exec.Command(xselArgs[0], xselArgs[1:]...)

	// pass the input string to the standard input of xsel
	clipCmd.Stdin = bytes.NewBuffer([]byte(input))

	// run the command
	err := clipCmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("copied to clipboard.\n")
}

// Returns a pointer to a Bookmark if it can be found within a list of bookmarks
func getBmark(bmarks Bookmarks, user string, site string) (*Bookmark, error) {

	// for each bookmark in the given bmarks...
	for _, bmark := range bmarks {

		// the identifier is 'user@site'
		id := fmt.Sprintf("%s@%s", user, site)
		bmarkId := fmt.Sprintf("%s@%s", bmark.Username, bmark.Url)

		// if the given id matches the current bookmark's,
		// then return a pointer to it
		if id == bmarkId {
			return &bmark, nil
		}
	}

	return new(Bookmark), errors.New("bookmark could not be found")
}

// Load bookmarks from a .JSON backup file
func loadBookmarks(bookmarksFile string) (bmarks Bookmarks, err error) {

	// open the file
	jsonFile, err := os.Open(bookmarksFile)
	if err != nil {
		return bmarks, err
	}
	defer jsonFile.Close()

	// convert the file to a byte array
	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return bmarks, err
	}

	// encode the JSON to the Bookmarks struct array
	err = json.Unmarshal(jsonBytes, &bmarks)
	if err != nil {
		return bmarks, err
	}

	return bmarks, nil
}

// Main program logic
func main() {

	// file to open
	var bookmarksFile string

	// the default bookmarks file location is $HOME/.CryptopassBookmarks.txt:
	defaultFile := fmt.Sprintf("%s/.CryptopassBookmarks.txt", os.Getenv("HOME"))

	// flag to enable custom path to bookmarks file...
	flag.StringVar(&bookmarksFile, "b", defaultFile, "bookmarks file")
	flag.Usage = usage
	flag.Parse()

	// load the bookmarks file...
	bmarks, err := loadBookmarks(bookmarksFile)
	if err != nil {
		log.Fatal(err)
	}

	// number of arguments remaining after flags are parsed
	narg := len(os.Args) - flag.NArg()

	if narg < 1 {
		fmt.Printf("insufficient arguments given\n")
		usage()
		return
	}

	// command parsing
	switch os.Args[narg] {
	case "help":
		// print usage
		usage()
	case "ls":
		// list bookmarks
		list(bmarks)
	case "find":
		// filter the bookmarks, then list them
		bmarks = filterList(bmarks, os.Args[narg+1])
		list(bmarks)
	case "open":
		// the bookmark given by the user
		givenBmark := os.Args[narg+1]

		// get the position of the last '@'
		i := strings.LastIndex(givenBmark, "@")

		user := givenBmark[:i]   // user is everything before the last '@'
		site := givenBmark[i+1:] // url is everything after the last '@' (not including it)

		fmt.Printf("user:%s\nsite:%s\n", user, site)

		// get pointer to Bookmark that matches the given user+site
		bmark, err := getBmark(bmarks, user, site)
		if err != nil {
			log.Fatal(err)
		}

		// user input's secret...
		secret, err := inputSecret()
		if err != nil {
			log.Fatal(err)
		}

		// generate the password from the given secret
		password := genPassword(secret, bmark)

		// copy the password to the clipboard
		clipboard(password)
	default:
		fmt.Printf("unknown command `%s'\n", os.Args[narg])
		usage()
	}

	return
}
