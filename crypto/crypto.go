// cpass/crypto
// this package defines crypto helpers. including password generation via
// PBKDF2 and age encryption helper functions

package crypto

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations int = 5000
	keyLen     int = 32
)

// generate a password
func CryptoPass(secret []byte, salt []byte, length int) string {
	var (
		pbkdf2Hmac []byte
		b64        string
	)

	pbkdf2Hmac = pbkdf2.Key(secret, salt, iterations, keyLen, sha256.New)
	b64 = base64.StdEncoding.EncodeToString(pbkdf2Hmac)

	return b64[:length]
}
