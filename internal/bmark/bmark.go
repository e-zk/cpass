package bmark

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/e-zk/cpass/internal/crypto"
)

type Bookmark struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Length   int    `json:"length"`
}
type Bookmarks []Bookmark

func LoadFromFile(filePath string) (bmarks Bookmarks, err error) {
	// open the file
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// convert the file to a byte array
	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	// encode the JSON to the Bookmarks struct array
	err = json.Unmarshal(jsonBytes, &bmarks)
	if err != nil {
		return nil, err
	}

	return bmarks, nil
}

func (b Bookmark) GenPassword(secret []byte) string {
	// the salt is "user@site"
	salt := fmt.Sprintf("%s@%s", b.Username, b.Url)
	return crypto.CryptoPass(secret, []byte(salt), b.Length)
}

func (bmarks Bookmarks) Get(user string, site string) (*Bookmark, error) {
	// for each bookmark in the given bmarks...
	for _, b := range bmarks {

		// the identifier is 'user@site'
		id := fmt.Sprintf("%s@%s", user, site)
		bmarkId := fmt.Sprintf("%s@%s", b.Username, b.Url)

		// if the given id matches the current bookmark's  then return a pointer
		if id == bmarkId {
			return &b, nil
		}
	}

	return new(Bookmark), errors.New("bookmark could not be found")
}

func (bmarks Bookmarks) String() (out string) {
	for _, b := range bmarks {
		out = fmt.Sprintf("%s%s@%s (%d)\n", out, b.Username, b.Url, b.Length)
	}
	return out
}

func (bmarks Bookmarks) Filter(filter string) (out Bookmarks) {
	for _, b := range bmarks {
		fullId := fmt.Sprintf("%s@%s", b.Username, b.Url)
		if strings.Contains(fullId, filter) {
			out = append(out, b)
		}
	}
	return out
}
