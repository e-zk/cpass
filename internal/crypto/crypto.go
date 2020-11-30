package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 5000
)

// Apply PBKDF2 with given secret and salt, output base64 encoded key
func CryptoPass(secret []byte, salt []byte, length int) string {
	// generate the pbkdf2 key based on the input values
	pbkdf2Hmac := pbkdf2.Key(secret, salt, iterations, 32, sha256.New)

	// encode the resulting pbkdf2 key in base64
	b64Encoded := base64.StdEncoding.EncodeToString([]byte(pbkdf2Hmac))

	// cut the encoded key down to the given length; this is the final password
	return b64Encoded[:length]
}
