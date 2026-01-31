package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func unzipString(str string) (string, error) {
	if str == "" {
		return "", errors.New("string is empty")
	}

	var builder strings.Builder

	runes := []rune(str)
	var prev rune
	hasPrev := false
	escape := false
	for _, r := range runes {
		if escape {
			prev = r
			hasPrev = true
			escape = false
			continue
		}
		if unicode.IsDigit(r) {
			if !hasPrev {
				return "", errors.New("invalid string")
			}
			n := int(r - '0')
			builder.WriteString(strings.Repeat(string(prev), n))
			hasPrev = false
			continue
		}
		if hasPrev {
			builder.WriteRune(prev)
			hasPrev = false
		}
		if r == '\\' {
			escape = true
			continue
		}
		prev = r
		hasPrev = true
	}
	if escape {
		return "", errors.New("escape subsequence at the end")
	}
	if hasPrev {
		builder.WriteRune(prev)
	}
	return builder.String(), nil
}
func main() {
	s1 := "a12ff\\34"
	result, err := unzipString(s1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unzipString error: %v", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
