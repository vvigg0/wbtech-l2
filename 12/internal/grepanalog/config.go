package grepanalog

// Config описывает параметры фильтрации и форматирования вывода,
// соответствующие флагам утилиты grep.
type Config struct {
	After, Before int
	CountOnly     bool
	IgnoreCase    bool
	Invert        bool
	Fixed         bool
	LineNum       bool
	Pattern       string
}
