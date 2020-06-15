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
	"runtime"
	"strings"

	// external dependencies...
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	iterations    = 5000                               // pbkdf2 iterations
	xselPath      = "xsel"                             // path to xsel(1)
	wslClipPath   = "/mnt/c/Windows/system32/clip.exe" // windows clip.exe path
	winClipPath   = "C:\\Windows\\system32\\clip.exe"  // windows clip.exe path
	osReleasePath = "/proc/sys/kernel/osrelease"
	printWarn     = "WARNING: will print password to stdout\n"
	secretPrompt  = "secret (will not echo): " // prompt for secret
)

var (
	xselArgs = []string{xselPath, "-i"} // xsel(1) arguments
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

	fmt.Printf("usage: %s [-p] [-b path] command [args]\n\n", os.Args[0])
	fmt.Printf("where:\n")
	fmt.Printf("\t-b path\t\tpath to bookmarks file\n")
	fmt.Printf("\t-p\t\tprint the password to stdout instead of piping to clipboard command\n")
	fmt.Printf("\n")
	fmt.Printf("valid commands:\n")
	fmt.Printf("\thelp\t\t\tprint this help message\n")
	fmt.Printf("\tls\t\t\tlist available bookmarks\n")
	fmt.Printf("\tfind 'string'\t\tsearch for a password containing a string\n")
	fmt.Printf("\topen user@site\topen bookmark with id 'user@site'\n")
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

	// use terminal API to read user password without echoing
	secret, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")

	// return the user's secret
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

// Checks which OS we are running on;
// if it is WSL it returns "WSL"
func getOS() string {

	// get runtime GOOS
	ret := runtime.GOOS

	// check if we are running on WSL by examining version string
	// located in /proc
	if ret == "linux" {
		verBytes, err := ioutil.ReadFile(osReleasePath)
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasSuffix(string(verBytes), "Microsoft\n") {
			ret = "WSL"
		}
	}

	// return OS
	return ret
}

// Copies a string to the clipboard via xsel(1)
func clipboard(input string) error {

	// command variable
	var clipCmd *exec.Cmd

	// if we are running on WSL or Windows, use windows' clip.exe
	if getOS() == "WSL" {
		clipCmd = exec.Command(wslClipPath)
	} else if getOS() == "windows" {
		clipCmd = exec.Command(winClipPath)
	} else {
		// xsel command
		clipCmd = exec.Command(xselArgs[0], xselArgs[1:]...)
	}

	// pass the input string to the standard input of the clipboard command
	clipCmd.Stdin = bytes.NewBuffer([]byte(input))

	// run the command
	err := clipCmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("copied to clipboard.\n")
	return nil
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

	// test if arguments are insufficient...
	if len(os.Args) == 1 {
		fmt.Printf("insufficient arguments given\n")
		usage()
		return
	}

	// file to open
	var bookmarksFile string
	var printPasswd bool

	// the default bookmarks file location is $HOME/.CryptopassBookmarks.txt:
	defaultFile := fmt.Sprintf("%s/.CryptopassBookmarks.txt", os.Getenv("HOME"))

	//
	defaultPrint := false

	// flag to enable custom path to bookmarks file...
	flag.StringVar(&bookmarksFile, "b", defaultFile, "bookmarks file")
	flag.BoolVar(&printPasswd, "p", defaultPrint, "print password to stdout, instead of passing it to xclip")
	flag.Usage = usage // enable custom usage function
	flag.Parse()       // parse flags

	// number of arguments remaining after flags are parsed
	narg := len(os.Args) - flag.NArg()

	// if the number of arguments remaining are less than one, fail and return
	// usage information
	if narg < 1 {
		fmt.Printf("insufficient arguments given\n")
		usage()
		return
	}

	// load the bookmarks file...
	bmarks, err := loadBookmarks(bookmarksFile)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// command parsing
	switch os.Args[narg] {
	case "help":
		usage()
	case "ls":
		list(bmarks)
	case "find":
		// filter the bookmarks, then list them
		bmarks = filterList(bmarks, os.Args[narg+1])
		list(bmarks)
	case "open":

		// print warning message if applicable
		if printPasswd {
			fmt.Printf(printWarn)
		}

		// the bookmark given by the user
		givenBmark := os.Args[narg+1]

		// get the position of the last '@'
		i := strings.LastIndex(givenBmark, "@")

		user := givenBmark[:i]   // user is everything before the last '@'
		site := givenBmark[i+1:] // url is everything after the last '@' (not including it)

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

		// print the password to stdout if -s is set;
		// if not set, then copy the password to the clipboard via xsel
		if printPasswd {
			fmt.Printf("%s\n", password)
		} else {
			// copy the password to the clipboard
			err = clipboard(password)
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		// if any other command is given, show an error message and usage information
		fmt.Printf("unknown command `%s'\n", os.Args[narg])
		usage()
	}

	return
}
