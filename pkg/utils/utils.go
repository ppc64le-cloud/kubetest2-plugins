package utils

import (
	"math/rand"
	"time"
)

// aASCII - the ASCII value of letter 'a'
const aASCII = 97

// RandString returns a string of lower-case alphabets of required length.
// The random string generated depends on seed, and is not cryptographically secure.
func RandString(length int) string {
	var genString string
	rand.Seed(time.Now().UnixNano())
	for ; length > 0; length-- {
		randChar := rand.Intn(26) + aASCII
		genString += string(rune(randChar))
	}
	return genString
}
