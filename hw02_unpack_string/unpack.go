package hw02unpackstring

import (
	"errors"
	"fmt"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	None   = 0
	Digit  = iota
	Letter = iota
	Shield = iota
)

func checkType(s rune) int {
	switch {
	case s == '\\':
		return Shield
	case unicode.IsDigit(s):
		return Digit
	default:
		return Letter
	}
}

func Unpack(s string) (string, error) {
	// Place your code here.

	if len(s) == 0 {
		return "", nil
	}

	input := []rune(s)
	output := make([]rune, 0, len(s)*2)
	first := input[0]
	lastType := None

	fmt.Println(string(input))

	switch checkType(first) {
	case Digit:
		return "", ErrInvalidString
	case Letter:
		lastType = Letter
		output = append(output, first)
	case Shield:
		lastType = Shield
	}

	for _, r := range input[1:] {
		switch {
		case lastType == Digit && checkType(r) == Digit:
			return "", ErrInvalidString
		case lastType == Shield:
			lastType = Letter
			output = append(output, r)
		case lastType == Letter && checkType(r) == Shield:
			lastType = Shield
		case lastType == Letter && r == '0':
			output = output[:len(output)-1]
			lastType = Digit
		case lastType == Letter && checkType(r) == Digit:
			for i := 0; i < int(r-'0')-1; i++ {
				output = append(output, output[len(output)-1])
			}
			lastType = Digit
		default:
			output = append(output, r)
			lastType = Letter
		}
	}

	return string(output), nil
}
