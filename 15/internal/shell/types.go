package shell

type tokenType int

// Типы токенов, используемые лексическим анализатором shell.
const (
	// WORD представляет обычное слово (имя команды, аргумент или путь).
	WORD tokenType = iota

	// PIPE представляет оператор конвейера "|".
	PIPE

	// AND представляет логический оператор "&&".
	AND

	// OR представляет логический оператор "||".
	OR

	// IN представляет оператор перенаправления ввода "<".
	IN

	// OUT представляет оператор перенаправления вывода ">".
	OUT
)

type token struct {
	ttype tokenType
	value string
}

type op int

const (
	opAnd op = iota
	opOr
)

func (o op) String() string {
	switch o {
	case opAnd:
		return "&&"
	case opOr:
		return "||"
	default:
		return "<?>"
	}
}

type expression struct {
	jobs []job
	ops  []op
}

type job struct {
	pipeline []command
}

type command struct {
	argv []string
	in   string
	out  string
}
