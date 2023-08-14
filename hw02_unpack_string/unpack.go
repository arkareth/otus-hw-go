package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	rs := []rune(s)
	sb := strings.Builder{}

	var skip int
	var escape bool
	for k, v := range rs {
		if skip > 0 {
			skip--
			continue
		}

		switch {
		case escape:
			escape = false
			if !unicode.IsDigit(v) && v != '\\' {
				return "", ErrInvalidString
			}
			if len(rs) >= k+2 && unicode.IsDigit(rs[k+1]) {
				skip = 1
				if rs[k+1] == 0 {
					continue
				}
				i, _ := strconv.Atoi(string(rs[k+1]))
				sb.WriteString(strings.Repeat(string(v), i))
				continue
			}
			sb.WriteRune(v)
			continue
		case unicode.IsDigit(v):
			if k == 0 {
				return "", ErrInvalidString
			}
			if unicode.IsDigit(rs[k-1]) {
				return "", ErrInvalidString
			}
			i, _ := strconv.Atoi(string(v))
			sb.WriteString(strings.Repeat(string(rs[k-1]), i-1))
			continue
		case v == '\\':
			escape = true
			continue
		case len(rs) >= k+2 && rs[k+1] == '0':
			skip = 1
			continue
		default:
			sb.WriteRune(v)
			continue
		}
	}

	return sb.String(), nil
}
