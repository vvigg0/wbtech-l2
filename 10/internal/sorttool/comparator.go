package sorttool

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// Comparator сравнивает две строки и возвращает
// -1 ,если a<b ;
// ноль, если a==b ;
// 1, если a>b ;
// в соответствии с правилами сортировки
type Comparator func(a, b string) int

// BuildComparator создает компаратор строк на основе Config.
// Компаратор выполняет извлечение столбца, необязательную нормализацию
// и сравнение в одном из заданных режимов: числовое, по месяцам, в удобочитаемом формате
// или лексикографическом(по умолчанию).
// В случае указания несовместимых параметров возвращается ошибка.
func BuildComparator(cfg *Config) (Comparator, error) {
	modes := 0
	if cfg.Numeric {
		modes++
	}
	if cfg.Month {
		modes++
	}
	if cfg.Human {
		modes++
	}
	if modes > 1 {
		return nil, errors.New("allowed only one of -n, -M, -h")
	}
	normalize := func(s string) string { return s }
	if cfg.TrimSpace {
		normalize = func(s string) string {
			return strings.TrimRightFunc(s, unicode.IsSpace)
		}
	}
	keyFn := func(line string) string { return normalize(line) }
	if cfg.KeyColumn > 0 {
		col := cfg.KeyColumn
		keyFn = func(line string) string {
			line = normalize(line)
			fields := strings.Split(line, cfg.Separator)
			if col-1 < 0 || col-1 >= len(fields) {
				return ""
			}
			return fields[col-1]
		}
	}
	var keyCmp func(ka, kb string) int = func(ka, kb string) int {
		if ka < kb {
			return -1
		}
		if ka > kb {
			return 1
		}
		return 0
	}

	switch {
	case cfg.Numeric:
		keyCmp = func(ka, kb string) int {
			na := parseInt(ka)
			nb := parseInt(kb)
			if na < nb {
				return -1
			}
			if na > nb {
				return 1
			}
			return 0
		}
	case cfg.Month:
		keyCmp = func(ka, kb string) int {
			ma := monthIndex(ka)
			mb := monthIndex(kb)
			if ma < mb {
				return -1
			}
			if ma > mb {
				return 1
			}
			return 0
		}
	case cfg.Human:
		keyCmp = func(ka, kb string) int {
			ha := parseHuman(ka)
			hb := parseHuman(kb)
			if ha < hb {
				return -1
			}
			if ha > hb {
				return 1
			}
			return 0
		}
	}
	cmp := func(a, b string) int {
		ka := keyFn(a)
		kb := keyFn(b)

		if c := keyCmp(ka, kb); c != 0 {
			if cfg.Reverse {
				return -c
			}
			return c
		}
		a = normalize(a)
		b = normalize(b)
		if a < b {
			if cfg.Reverse {
				return 1
			}
			return -1
		}
		if a > b {
			if cfg.Reverse {
				return -1
			}
			return 1
		}
		return 0
	}
	return cmp, nil
}
func parseInt(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(f)
}
func monthIndex(s string) int {
	s = strings.TrimSpace(s)
	if len(s) < 3 {
		return 0
	}
	m := strings.ToLower(s[:3])
	switch m {
	case "jan":
		return 1
	case "feb":
		return 2
	case "mar":
		return 3
	case "apr":
		return 4
	case "may":
		return 5
	case "jun":
		return 6
	case "jul":
		return 7
	case "aug":
		return 8
	case "sep":
		return 9
	case "oct":
		return 10
	case "nov":
		return 11
	case "dec":
		return 12
	default:
		return 13
	}
}
func parseHuman(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	last := s[len(s)-1]
	mult := int64(1)
	num := s[:len(s)-1]
	switch last {
	case 'K', 'k':
		mult = 1024
	case 'M', 'm':
		mult = 1024 * 1024
	case 'G', 'g':
		mult = 1024 * 1024 * 1024
	case 'T', 't':
		mult = 1024 * 1024 * 1024 * 1024
	}
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0
	}
	return int64(f * float64(mult))
}
