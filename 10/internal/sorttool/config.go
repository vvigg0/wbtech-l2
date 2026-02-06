package sorttool

// Config описывает параметры сортировки
type Config struct {
	KeyColumn int
	Numeric   bool
	Reverse   bool
	Unique    bool
	Month     bool
	TrimSpace bool
	Check     bool
	Human     bool
	Separator string
}
