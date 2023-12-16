package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	Digit  = 1
	Leter  = iota
	Shield = iota
)

func Unpack(s string) (string, error) {
	// Place your code here.

	if len(s) == 0 {
		return s, nil
	}

	input := []rune(s)
	if unicode.IsDigit(input[0]) {
		return "", ErrInvalidString
	}

	output := make([]rune, 0, len(s))
	lastRune := input[0]
	output = append(output, lastRune)

	if len(input) > 1 {
		for _, r := range input[1:] {
			switch {
			case unicode.IsDigit(lastRune) && unicode.IsDigit(r):
				return "", ErrInvalidString
			case unicode.IsLetter(lastRune) && r == '0':
				output = output[:len(output)-1]
			case unicode.IsLetter(lastRune) && unicode.IsDigit(r):
				for i := 0; i < int(r-'0')-1; i++ {
					output = append(output, lastRune)
				}
				lastRune = r
			case unicode.IsLetter(r):
				output = append(output, r)
				lastRune = r
			}
		}
	}
	return string(output), nil
}
