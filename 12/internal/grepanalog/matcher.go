package grepanalog

import (
	"regexp"
	"strings"
)

// BuildMatcher создаёт и возвращает функцию сопоставления строки с шаблоном
// на основе переданной конфигурации.
//
// Возвращаемая функция принимает строку входного потока и возвращает true,
// если строка считается совпавшей с шаблоном с учётом всех флагов
// (IgnoreCase, Fixed, Invert).
//
// В случае ошибки компиляции регулярного выражения возвращается ошибка.
func BuildMatcher(cfg *Config) (func(string) bool, error) {
	pattern := cfg.Pattern
	if cfg.IgnoreCase {
		pattern = strings.ToLower(cfg.Pattern)
	}

	var reg *regexp.Regexp
	var err error
	if !cfg.Fixed {
		reg, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
	}
	match := func(line string) bool {
		if cfg.IgnoreCase {
			line = strings.ToLower(line)
		}

		var res bool
		if cfg.Fixed {
			res = strings.Contains(line, pattern)
		} else {
			res = reg.MatchString(line)
		}

		if cfg.Invert {
			return !res
		}
		return res
	}
	return match, nil
}
