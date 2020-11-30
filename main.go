package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/e-zk/cpass/internal/bmark"

	"github.com/atotto/clipboard"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	wslClipPath   = "/mnt/c/Windows/system32/clip.exe"
	winClipPath   = "C:\\Windows\\system32\\clip.exe"
	osReleasePath = "/proc/sys/kernel/osrelease" // TODO change to LINUX releasePath
	printWarn     = "WARNING: will print password to stdout\n"
	secretPrompt  = "secret (will not echo): "
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

		if strings.HasSuffix(string(verBytes), "microsoft-standard\n") {
			ret = "WSL"
		}
	}

	// return OS
	return ret
}

//// Copies a string to the clipboard via xsel(1)
//func clipboard(input string) error {
//
//	// command variable
//	var clipCmd *exec.Cmd
//
//	// if we are running on WSL or Windows, use windows' clip.exe
//	if getOS() == "WSL" {
//		clipCmd = exec.Command(wslClipPath)
//	} else if getOS() == "windows" {
//		clipCmd = exec.Command(winClipPath)
//	} else {
//		clipCmd = exec.Command(xselArgs[0], xselArgs[1:]...)
//	}
//
//	// pass the input string to the standard input of the clipboard command
//	clipCmd.Stdin = bytes.NewBuffer([]byte(input))
//
//	// run the command
//	err := clipCmd.Run()
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("copied to clipboard.\n")
//	return nil
//}

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

		// get the position of the last '@'
		i := strings.LastIndex(givenBmark, "@")

		// extract user + site
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

		// print the password to stdout if -s is set;
		// if not set, then copy the password to the clipboard
		if printPasswd {
			fmt.Printf("%s\n", password)
		} else {
			err = clipboard.WriteAll(password)
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
