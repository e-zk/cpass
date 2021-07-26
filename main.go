package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/e-zk/cpass/store"
	"github.com/e-zk/cpass/term"

	"github.com/atotto/clipboard"
	"github.com/e-zk/subc"
	"github.com/e-zk/wslcheck"
)

const (
	wslClipPath = "/mnt/c/Windows/system32/clip.exe"
	warnPrint   = "warning: will print password to standard output"
	defaultLen  = 16
)

var (
	err         error
	storePath   string
	printPasswd bool
)

// wrapper for Fprintf to print to stdout
func errPrint(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// main usage/help message
func usage() {
	errPrint("usage: cpass [command] [args]\n\n")
	errPrint("where [command] can be:\n")
	errPrint("  help  show this help message\n")
	errPrint("  ls    list password entries\n")
	errPrint("  open  open/view a password entry\n")
	errPrint("  save  save/add a new password entry\n")
	errPrint("  rm    remove password entry\n")
	errPrint("\n")
	errPrint("for help with subcommands type: cpass [command] -h\n")
}

//
func clip(text string) (err error) {
	wsl, err := wslcheck.Check()
	if err != nil {
		return err
	}

	if wsl {
		cmd := exec.Command(wslClipPath)
		cmd.Stdin = bytes.NewBuffer([]byte(text))
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		err = clipboard.WriteAll(text)
		if err != nil {
			return err
		}
	}

	return nil
}

// list all password entries
func list() {
	//help := subc.Sub("ls").Arg(0)
	//if help == "help" { // subc might do this for us
	//	subc.Sub("ls").Usage()
	//	return
	//}

	s, err := store.NewStore(storePath)
	if err != nil {
		fmt.Println(s)
		log.Fatal(err)
	}

	entries, err := s.Entries()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", entries.String())
}

// save a new password entry
func save(printPasswd bool, entryLen int) {
	entryId := subc.Sub("save").Arg(0)
	//if entryId == "help" {
	//	subc.Sub("save").Usage()
	//	return
	//}

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
func remove(force bool) {
	entryId := subc.Sub("rm").Arg(0)

	// create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		log.Fatal(err)
	}

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
func open(printPasswd bool) {
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

	entryId := subc.Sub("open").Arg(0)
	//if entryId == "help" {
	//	openFlag.Usage()
	//	return
	//}

	e := entries.Get(entryId)
	if e == nil {
		log.Fatal("bookmark not found")
	}

	// ask for secret
	secret, err := term.AskPasswd()
	if err != nil {
		log.Fatal(err)
	}

	genPasswd := e.GenPassword(secret)

	if printPasswd {
		fmt.Printf("%s\n", genPasswd)
		return
	}

	err = clip(genPasswd)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(0 | log.Lshortfile)
	log.SetPrefix("cpass: ")

	// get default password store location
	configHome, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	defaultStore := configHome + "/cpass/bookmarks.json"
	storePath = defaultStore

	var (
		entryLen    int
		printPasswd bool
		force       bool
	)

	subc.Usage = usage

	subc.Sub("help")

	subc.Sub("ls").StringVar(&storePath, "s", defaultStore, "path to password store")
	subc.Sub("ls").Usage = func() {
		errPrint("list all entries in store\n")
		errPrint("usage: cpass ls [-s store]\n\n")
		errPrint("  -s store  use given password store\n")
	}

	subc.Sub("save").StringVar(&storePath, "s", defaultStore, "path to password store")
	subc.Sub("save").BoolVar(&printPasswd, "p", false, "print password to stdout")
	subc.Sub("save").IntVar(&entryLen, "l", defaultLen, "specify password length")
	subc.Sub("save").Usage = func() {
		errPrint("save a new password entry\n")
		errPrint("usage: cpass save [-l len] [-p] [-s store] <user@site>\n\n")
		errPrint("  -l len    specify password length (default 16)\n")
		errPrint("  -p        print password to stdout\n")
		errPrint("  -s store  use password store\n")
	}

	subc.Sub("rm").StringVar(&storePath, "s", defaultStore, "path to password store")
	subc.Sub("rm").BoolVar(&force, "f", false, "force remove entry (do not prompt)")
	subc.Sub("rm").Usage = func() {
		errPrint("remove a password entry\n")
		errPrint("usage: cpass rm [-f] [-s store] <user@site>\n\n")
		errPrint("    -f          force - do not prompt before removing")
		errPrint("    -s store    use password store")
	}

	subc.Sub("open").StringVar(&storePath, "s", defaultStore, "path to password store")
	subc.Sub("open").BoolVar(&printPasswd, "p", false, "print password to stdout")
	subc.Sub("open").Usage = func() {
		errPrint("open a password entry\n")
		errPrint("usage: cpass open [-p] [-s store] [-k key_file] <user@site>\n\n")
		errPrint("  -p        print password to stdout\n")
		errPrint("  -s store  use password store\n")
		//errPrint("    -k key_file    supply key_file when using an encrypted store\n")
	}

	subcommand, err := subc.Parse()
	if err != nil {
		log.Fatal(err)
	}

	switch subcommand {
	case "help":
		usage()
	case "ls":
		list()
	case "save":
		save(printPasswd, entryLen)
	case "rm":
		remove(force)
	case "open":
		open(printPasswd)
	default:
		errPrint("unknown command `%s'\n", os.Args[1])
		usage()
		os.Exit(1)
	}

}
