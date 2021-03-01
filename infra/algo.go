package infra

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// Hmac Generate a hex hash value with the key,
// expects: hmacmd5, hmacsha1, hmacsha224, hmacsha256, hmacsha384, hmacsha512.
func Hmac(method, key, val string) string {
	var f func() hash.Hash

	switch method {
	case "hmacmd5":
		f = md5.New
	case "hmacsha1":
		f = sha1.New
	case "hmacsha224":
		f = sha256.New224
	case "hmacsha256":
		f = sha256.New
	case "hmacsha384":
		f = sha512.New384
	case "hmacsha512":
		f = sha512.New
	default:
		return val
	}
	h := hmac.New(f, []byte(key))
	h.Write([]byte(val)) // nolint: errCheck
	return hex.EncodeToString(h.Sum(nil))
}
