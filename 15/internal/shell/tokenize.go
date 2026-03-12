package shell

import (
	"errors"
	"strings"
	"unicode"
)

func tokenize(s string) ([]token, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	const (
		qNone = iota
		qSingle
		qDouble
	)

	q := qNone
	esc := false

	var out []token
	var b strings.Builder

	flushWord := func() {
		if b.Len() == 0 {
			return
		}
		out = append(out, token{ttype: WORD, value: b.String()})
		b.Reset()
	}

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if esc {
			b.WriteByte(ch)
			esc = false
			continue
		}

		if q == qNone && ch == '\\' {
			esc = true
			continue
		}

		if q == qSingle {
			if ch == '\'' {
				q = qNone
				continue
			}
			b.WriteByte(ch)
			continue
		}

		if q == qDouble {
			if ch == '"' {
				q = qNone
				continue
			}
			b.WriteByte(ch)
			continue
		}

		if unicode.IsSpace(rune(ch)) {
			flushWord()
			continue
		}

		switch ch {
		case '\'':
			q = qSingle
			continue
		case '"':
			q = qDouble
			continue
		case '|':
			flushWord()
			if i+1 < len(s) && s[i+1] == '|' {
				out = append(out, token{ttype: OR})
				i++
			} else {
				out = append(out, token{ttype: PIPE})
			}
			continue
		case '&':
			flushWord()
			if i+1 < len(s) && s[i+1] == '&' {
				out = append(out, token{ttype: AND})
				i++
				continue
			}
			return nil, errors.New("одиночный &")
		case '<':
			flushWord()
			out = append(out, token{ttype: IN})
			continue
		case '>':
			flushWord()
			out = append(out, token{ttype: OUT})
			continue
		default:
			b.WriteByte(ch)
		}
	}

	if esc {
		return nil, errors.New("dangling escape")
	}
	if q != qNone {
		return nil, errors.New("незакрытая кавычка")
	}

	flushWord()
	return out, nil
}
