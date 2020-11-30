package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/e-zk/cpass/internal/bmark"

	"github.com/atotto/clipboard"
	"github.com/e-zk/wslcheck"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	wslClipPath  = "/mnt/c/Windows/system32/clip.exe"
	printWarn    = "WARNING: will print password to stdout\n"
	secretPrompt = "secret (will not echo): "
)

// Prints program usage information
func usage() {

	fmt.Printf("usage:\n")
	fmt.Printf("    %s [flags] <subcommand [args]>\n\n", os.Args[0])
	fmt.Printf("flags:\n")
	fmt.Printf("    -b path          path to bookmarks file\n")
	fmt.Printf("    -p               print the password to stdout instead of piping to clipboard\n\n")
	fmt.Printf("subcommands:\n")
	fmt.Printf("    help             print this help message\n")
	fmt.Printf("    ls               list available bookmarks\n")
	fmt.Printf("    find <string>    search for a password containing a string\n")
	fmt.Printf("    open <id>        open bookmark with id 'user@site'\n")
}

// Get secret from user
func inputSecret() ([]byte, error) {

	fmt.Printf(secretPrompt)

	// use terminal API to read user password without echoing
	secret, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")

	return secret, err
}

/* Send string to clipboard */
func clip(input string) error {
	var (
		wsl bool
		err error
	)

	wsl, err = wslcheck.Check()
	if err != nil {
		return err
	}

	// if we are running on WSL; execute clip.exe
	// otherwise use clipboard module
	if wsl {
		cmd := exec.Command(wslClipPath)
		cmd.Stdin = bytes.NewBuffer([]byte(input))
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		err = clipboard.WriteAll(input)
		if err != nil {
			return err
		}
	}

	return nil
}

// Main program logic
func main() {

	if len(os.Args) == 1 {
		fmt.Printf("insufficient arguments given\n")
		usage()
		return
	}

	var (
		defaultFile   string
		bookmarksFile string
		printPasswd   bool
		narg          int
		defaultPrint  bool = false
	)

	// get config dir
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		defaultFile = fmt.Sprintf("%s/.config/cpass/bookmarks.json", os.Getenv("HOME"))
	} else {
		defaultFile = fmt.Sprintf("%s/cpass/bookmarks.json", configHome)
	}

	// flags
	flag.StringVar(&bookmarksFile, "b", defaultFile, "")
	flag.BoolVar(&printPasswd, "p", defaultPrint, "")
	flag.Usage = usage
	flag.Parse()

	// if the number of arguments remaining are less than one, fail and return
	// usage information
	if narg = len(os.Args) - flag.NArg(); narg < 1 {
		fmt.Printf("insufficient arguments given\n")
		usage()
		return
	}

	// load the bookmarks file...
	bmarks, err := bmark.LoadFromFile(bookmarksFile)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	// subcommand parsing
	switch os.Args[narg] {
	case "help":
		usage()
	case "ls":
		print(bmarks.String())
	case "find":
		// filter the bookmarks, then list them
		bmarks = bmarks.Filter(os.Args[narg+1])
		print(bmarks.String())
	case "open":

		// print warning message if applicable
		if printPasswd {
			fmt.Printf(printWarn)
		}

		// the bookmark given by the user
		givenBmark := os.Args[narg+1]

		// extract user + site from 'user@site'
		// done this way so that everything after the first @ is the site
		i := strings.LastIndex(givenBmark, "@")
		user := givenBmark[:i]
		site := givenBmark[i+1:]

		// get Bookmark that matches the given user+site
		bmark, err := bmarks.Get(user, site)
		if err != nil {
			log.Fatal(err)
		}

		// user input's secret...
		secret, err := inputSecret()
		if err != nil {
			log.Fatal(err)
		}

		// generate the password from the given secret
		password := bmark.GenPassword(secret)

		// print the password to stdout if -p is set;
		// if not set, then copy the password to the clipboard
		if printPasswd {
			fmt.Printf("%s\n", password)
		} else {
			err = clip(password)
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		fmt.Printf("unknown command `%s'\n", os.Args[narg])
		usage()
		os.Exit(1)
	}
}
