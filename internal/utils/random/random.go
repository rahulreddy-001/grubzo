package random

import (
	"crypto/rand"
	"encoding/base64"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

func SecureAlphaNumeric(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := 0; i < n; {
		idx := int(b[i] & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i++
		} else {
			if _, err := rand.Read(b[i : i+1]); err != nil {
				panic(err)
			}
		}
	}
	return *(*string)(unsafe.Pointer(&b))
}

func HashPassword(pass string, salt string) string {
	hash := argon2.IDKey([]byte(pass), []byte(salt), 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash[:])
}
