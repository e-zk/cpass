package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	//"strings"

	"github.com/e-zk/cpass/store"
)

const (
	warnPrint = "warning: will print password to standard output"
)

// wrapper for Fprintf to print to stdout
func errPrint(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// subcommand flags
var (
	lsFlag   = flag.NewFlagSet("ls", flag.ExitOnError)
	findFlag = flag.NewFlagSet("find", flag.ExitOnError)
	openFlag = flag.NewFlagSet("open", flag.ExitOnError)
	saveFlag = flag.NewFlagSet("save", flag.ExitOnError)
	rmFlag   = flag.NewFlagSet("rm", flag.ExitOnError)
)

// give all flags help messages
func commandHelp() {
	lsFlag.Usage = func() {
		errPrint("list all entries in store")
		errPrint("usage: cpass ls [-s store]\n\n")
		errPrint("    -s store    use given password store\n")
	}

	findFlag.Usage = func() {
		errPrint("find password entries containing substring\n")
		errPrint("usage: cpass find [-s store] <substring>\n\n")
		errPrint("    -s store    use password store\n")
	}

	openFlag.Usage = func() {
		errPrint("open a password entry\n")
		errPrint("usage: cpass open [-p] [-s store] [-k key_file] <user@site>\n\n")
		errPrint("    -p             print password to stdout\n")
		errPrint("    -s store       use password store\n")
		errPrint("    -k key_file    supply key_file when using an encrypted store\n")
	}

	saveFlag.Usage = func() {
		fmt.Printf("usage: cpass save [-l len] [-p] [-s store] <user@site>\n\n")
	}

	rmFlag.Usage = func() {
		fmt.Printf("usage: cpass rm [-f] [-s store] <user@site>\n\n")
	}
}

// main usage/help message
func usage() {
	errPrint("usage: cpass [command] [args]\n\n")
	errPrint("where [command] can be:\n")
	errPrint("    help    show this help message\n")
	errPrint("    open    open/view a password entry\n")
	errPrint("    save    save/add a new password entry\n")
	errPrint("    rm      remove password entry\n")
	errPrint("\n")
	errPrint("for help with subcommands type: cpass [command] -h\n")
}

func main() {
	const (
		defaultPrint bool = false
		defaultLen   int  = 16
	)

	var (
		err          error
		configHome   string
		defaultStore string
		storePath    string
		printPasswd  bool
		force        bool
		length       int
		s            store.Store
	)

	log.SetFlags(0)
	log.SetPrefix("cpass: ")

	// get default password store location
	configHome, err = os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	defaultStore = configHome + "/cpass/bookmarks.json"

	// setup subcommand help messages
	commandHelp()

	// parse subcommand
	subcommand := os.Args[1]
	switch subcommand {
	case "help":
		usage()
	case "ls":
		lsFlag.StringVar(&storePath, "s", defaultStore, "")
		lsFlag.Parse(os.Args[2:])

		s, err := store.NewStore(storePath)
		if err != nil {
			log.Fatal(err)
		}

		entries, err := s.Entries()
		if err != nil {
			log.Fatal(err)
		}

		print(entries.String())

	case "save":
		saveFlag.IntVar(&length, "l", defaultLen, "")
		saveFlag.StringVar(&storePath, "s", defaultStore, "")
		saveFlag.BoolVar(&printPasswd, "p", defaultPrint, "")
		saveFlag.Parse(os.Args[2:])

		entryId := saveFlag.Arg(0)

		if entryId == "help" {
			saveFlag.Usage()
			return
		}

		s, err = store.NewStore(storePath)
		if err != nil {
			log.Fatal(err)
		}

		err = s.AddEntry(entryId, length)
		if err != nil {
			log.Fatal(err)
		}
	case "rm":
		rmFlag.StringVar(&storePath, "s", defaultStore, "")
		rmFlag.BoolVar(&force, "f", false, "")
		rmFlag.Parse(os.Args[2:])

		entryId := rmFlag.Arg(0)
		if entryId == "help" {
			rmFlag.Usage()
			return
		}

		s, err = store.NewStore(storePath)
		if err != nil {
			log.Fatal(err)
		}

		err = s.RemoveEntry(entryId)
		if err != nil {
			log.Fatal(err)
		}
	case "open":
		openFlag.StringVar(&storePath, "s", defaultStore, "")
		openFlag.BoolVar(&printPasswd, "p", defaultPrint, "")
		openFlag.Parse(os.Args[2:])

		if printPasswd {
			errPrint("%s\n", warnPrint)
		}

		s, err = store.NewStore(storePath)
		if err != nil {
			log.Fatal(err)
		}

		entries, err := s.Entries()
		if err != nil {
			log.Fatal(err)
		}

		entryId := openFlag.Arg(0)
		if entryId == "help" {
			openFlag.Usage()
			return
		}

		// this is kept here for future reference
		//i := strings.LastIndex(entryId, "@")
		//user := givenBark[:i]
		//url := entryId[i+1:]

		e := entries.Get(entryId)
		if e == nil {
			log.Fatal("bookmark not found")
		}

		// TODO : get secret input

		if printPasswd {
			fmt.Printf("%s\n", e.GenPassword([]byte("test")))
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown command `%s'\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}
