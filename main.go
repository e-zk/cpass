package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	//"strconv"
	"errors"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// iterations
	iterations = 5000

	// path to xsel(1)
	xselPath = "xsel"

	// password prompt
	enterSecretPrompt = "secret (will not echo): "
)

var (
	// xsel arguments
	xselArgs = []string{xselPath, "-i"}
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
	fmt.Printf("usage: %s <command> [<args>]\n\n", os.Args[0])
	fmt.Printf("command can be one of:\n")
	fmt.Printf("\thelp\t\t\tprint this help message\n")
	fmt.Printf("\tls\t\t\tlist available bookmarks\n")
	fmt.Printf("\tfind <string>\t\tfind a password containing <string>\n")
	fmt.Printf("\topen <user@site>\topen bookmark\n")
}

// Generate password from
func genPassword(secret []byte, user string, site string) string {

	salt := fmt.Sprintf("%s@%s", user, site)
	pbkdf2Hmac := pbkdf2.Key([]byte(secret), []byte(salt), 5000, 32, sha256.New)
	b64Encoded := base64.StdEncoding.EncodeToString([]byte(pbkdf2Hmac))

	return b64Encoded
}

// Get secret from user
func inputSecret() (string, error) {
	fmt.Printf(enterPasswdPrompt)
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")
	return password, err
}

// Generate and return a password using the given username, site and length
func getPass(username string, site string, length int) string {

	password, err := inputSecret()
	if err != nil {
		log.Fatal(err)
	}

	encodedDk := genPassword(password, username, site)
	return encodedDk[:length]
}

// Filter a list of bookmarks, based on a keyword and return bookmarks which match this keyword
func filterList(bmarks Bookmarks, filter string) Bookmarks {
	var outBmarks Bookmarks

	// if the filter string is within the Username or URL
	for _, bmark := range bmarks {
		if strings.Contains(bmark.Username, filter) {
			outBmarks = append(outBmarks, bmark)
		} else if strings.Contains(bmark.Url, filter) {
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

	// pass the password to the standard input of xsel
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
	for _, bmark := range bmarks {
		if bmark.Url == site && bmark.Username == user {
			return &bmark, nil
		}
	}
	return new(Bookmark), errors.New("bookmark could not be found")
}

// Main program logic
func main() {

	var bmarks Bookmarks

	testFileData := []byte(`
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
		`)

	// encode the JSON to the Bookmarks struct array
	jsonErr := json.Unmarshal(testFileData, &bmarks)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// argument parsing...
	switch os.Args[1] {
	case "help":
		// print usage
		usage()
	case "ls":
		// list bookmarks
		list(bmarks)
	case "find":
		// filter the bookmarks, then list them
		bmarks = filterList(bmarks, os.Args[2])
		list(bmarks)
	case "open":
		var user, site string

		// split the full password identifier (user@site) at '@'
		full := strings.Split(os.Args[2], "@")
		user = full[0]
		site = full[1]

		//
		bmark, err := getBmark(bmarks, user, site)
		if err != nil {
			log.Fatal(err)
		}

		// compute the password
		password := getPass(bmark.Username, bmark.Url, bmark.Length)

		// copy the password to the clipboard
		clipboard(password)

	}
}
