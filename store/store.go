// cpass/store
// this package describes the model of a password store as well as password
// entries

package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/e-zk/cpass/crypto"
)

const (
	jsonIndent = "  "
)

// Store errors
var (
	ErrStoreIsDir = errors.New("given path is a directory")
	ErrStoreExt   = errors.New("given store has incorrect extension")
	ErrStoreEnc   = errors.New("encrypted stores are not yet supported!")
)

// Entry errors
var (
	ErrEntryExists   = errors.New("entry already exists")
	ErrEntryNotExist = errors.New("entry does not exist")
)

// password store struct
type Store struct {
	FilePath  string
	StoreName string
	Encrypted bool
}

// password entry struct
type Entry struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Length   int    `json:"length"`
}

// typdef a slice of entry structs
type Entries []Entry

// create a new store struct based on the given password store filepath
func NewStore(path string) (Store, error) {
	var (
		s         os.FileInfo
		err       error
		base      string
		ext       string
		name      string
		encrypted bool
	)

	if s, err = os.Stat(path); err != nil {
		return Store{}, err
	} else if s.IsDir() {
		return Store{}, ErrStoreIsDir
	} else if ext = filepath.Ext(s.Name()); !(ext == ".json" || ext == ".age") {
		return Store{}, ErrStoreExt
	}

	base = s.Name()
	//ext = filepath.Ext(path)
	name = strings.TrimSuffix(base, ext)

	if ext == ".age" {
		encrypted = true
		name = strings.TrimSuffix(base, ".json")
	} else {
		encrypted = false
	}

	return Store{
		FilePath:  path,
		StoreName: name,
		Encrypted: encrypted,
	}, nil
}

// return list of bookmarks belonging to store
func (s Store) Entries() (es Entries, err error) {
	var (
		jsonBytes []byte
	)

	// TODO
	if s.Encrypted {
		return nil, ErrStoreEnc
	}

	// open the file
	jsonFile, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonBytes, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &es)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func (s Store) AddEntry(entryId string, length int) error {
	var (
		err       error
		username  string
		url       string
		jsonBytes []byte
		entry     Entry
		es        Entries
	)

	i := strings.LastIndex(entryId, "@")
	username = entryId[:i]
	url = entryId[i+1:]

	entry = Entry{
		Username: username,
		Url:      url,
		Length:   length,
	}

	// TODO
	if s.Encrypted {
		return ErrStoreEnc
	}

	// get entries
	es, err = s.Entries()
	if err != nil {
		return err
	}

	if es.Get(username+"@"+url) != nil {
		return ErrEntryExists
	}

	// append the new entry to the struct array
	es = append(es, entry)

	// marshal the data
	jsonBytes, err = json.MarshalIndent(es, "", jsonIndent)
	if err != nil {
		return err
	}

	// overwrite file
	err = ioutil.WriteFile(s.FilePath, jsonBytes, 640)
	if err != nil {
		return err
	}

	return nil
}

// delete a passwoard entry
func (s Store) RemoveEntry(entryId string) error {
	var (
		err       error
		index     int
		found     bool
		jsonBytes []byte
		es        Entries
	)

	// get entries
	es, err = s.Entries()
	if err != nil {
		return err
	}

	for i, e := range es {
		if e.Id() == entryId {
			index = i
			found = true
			break
		}
	}

	//
	if !found {
		return ErrEntryNotExist
	}

	// cut out the entry that is to be removed
	es = append(es[:index], es[index+1:]...)

	// marshal the new data
	jsonBytes, err = json.MarshalIndent(es, "", jsonIndent)
	if err != nil {
		return err
	}

	// overwrite file
	err = ioutil.WriteFile(s.FilePath, jsonBytes, 640)
	if err != nil {
		return err
	}

	return nil
}

// output string representation of a bookmark
func (e Entry) String() (out string) {
	out = fmt.Sprintf("%s@%s (%d)", e.Username, e.Url, e.Length)
	return out
}

// get entry id
func (e Entry) Id() string {
	return e.Username + "@" + e.Url
}

// generate password for password entry
func (e Entry) GenPassword(secret []byte) string {
	salt := fmt.Sprintf("%s@%s", e.Username, e.Url)
	return crypto.CryptoPass(secret, []byte(salt), e.Length)
}

// output string representation of a list of bookmarks
func (es Entries) String() (out string) {
	for _, e := range es {
		out = fmt.Sprintf("%s%s\n", out, e.String())
	}
	return out
}

// from a list of password entries, find the one matching the given username + url
func (es Entries) Get(givenId string) *Entry {
	//givenId := fmt.Sprintf("%s@%s", username, url)

	for _, e := range es {
		id := fmt.Sprintf("%s@%s", e.Username, e.Url)

		if id == givenId {
			return &e
		}
	}

	return nil
}
