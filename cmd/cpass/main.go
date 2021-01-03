package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/e-zk/cpass/store"
	"github.com/e-zk/cpass/term"
)

const (
	warnPrint  = "warning: will print password to standard output"
	defaultLen = 16
)

var (
	err          error
	defaultStore string
	storePath    string
	printPasswd  bool
)

// subcommand flags
var (
	lsFlag   = flag.NewFlagSet("ls", flag.ExitOnError)
	findFlag = flag.NewFlagSet("find", flag.ExitOnError)
	openFlag = flag.NewFlagSet("open", flag.ExitOnError)
	saveFlag = flag.NewFlagSet("save", flag.ExitOnError)
	rmFlag   = flag.NewFlagSet("rm", flag.ExitOnError)
)

// wrapper for Fprintf to print to stdout
func errPrint(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

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
		errPrint("save a new password entry\n")
		errPrint("usage: cpass save [-l len] [-p] [-s store] <user@site>\n\n")
		errPrint("    -l len      specify password length (default 16)\n")
		errPrint("    -p          print password to stdout\n")
		errPrint("    -s store    use password store\n")
	}

	rmFlag.Usage = func() {
		errPrint("usage: cpass rm [-f] [-s store] <user@site>\n\n")
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

// list all password entries
func list() {
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
}

// save a new password entry
func save() {
	var entryLen int

	saveFlag.IntVar(&entryLen, "l", defaultLen, "")
	saveFlag.StringVar(&storePath, "s", defaultStore, "")
	saveFlag.BoolVar(&printPasswd, "p", false, "")
	saveFlag.Parse(os.Args[2:])

	entryId := saveFlag.Arg(0)

	if entryId == "help" {
		saveFlag.Usage()
		return
	}

	s, err := store.NewStore(storePath)
	if err != nil {
		log.Fatal(err)
	}

	err = s.AddEntry(entryId, entryLen)
	if err != nil {
		log.Fatal(err)
	}
}

// remove a password entry
func remove() {
	var force bool

	rmFlag.StringVar(&storePath, "s", defaultStore, "")
	rmFlag.BoolVar(&force, "f", false, "")
	rmFlag.Parse(os.Args[2:])

	entryId := rmFlag.Arg(0)
	if entryId == "help" {
		rmFlag.Usage()
		return
	}

	// create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		log.Fatal(err)
	}

	// make sure the entry exists
	ok, err := s.EntryExists(entryId)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatalf("entry %s does not exist", entryId)
	}

	if !force {
		ch, err := term.Ask("remove entry " + entryId + "?")
		if err != nil {
			log.Fatal(err)
		}

		if !ch {
			errPrint("aborted.\n")
			return
		}
	}

	err = s.RemoveEntry(entryId)
	if err != nil {
		log.Fatal(err)
	}
}

// open a password entry
func open() {
	openFlag.StringVar(&storePath, "s", defaultStore, "")
	openFlag.BoolVar(&printPasswd, "p", false, "")
	openFlag.Parse(os.Args[2:])

	if printPasswd {
		errPrint("%s\n", warnPrint)
	}

	s, err := store.NewStore(storePath)
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

	e := entries.Get(entryId)
	if e == nil {
		log.Fatal("bookmark not found")
	}

	secret, err := term.AskPasswd()
	if err != nil {
		log.Fatal(err)
	}

	genPasswd := e.GenPassword(secret)

	if printPasswd {
		fmt.Printf("%s\n", genPasswd)
	}

	// TODO clipboard
}

func main() {
	// logging
	log.SetFlags(0)
	log.SetPrefix("cpass: ")

	// get default password store location
	configHome, err := os.UserConfigDir()
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
		list()
	case "save":
		save()
	case "rm":
		remove()
	case "open":
		open()
	default:
		errPrint("unknown command `%s'\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}
