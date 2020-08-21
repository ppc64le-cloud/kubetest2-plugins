package utils

import (
	"fmt"
	"github.com/lucasjones/reggen"
)

// RandString generates n number of random char string
func RandString(n int) (string, error) {
	return reggen.Generate(fmt.Sprintf("[a-z]{%d}", n), 2)
}

func GenerateBootstrapToken() (string, error) {
	return reggen.Generate("[a-z0-9]{6}\\.[a-z0-9]{16}", 1)
}
