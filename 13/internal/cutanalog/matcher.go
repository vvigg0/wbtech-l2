package cutanalog

import (
	"fmt"
	"strings"
)

// BuildMatcher возвращает функцию, которая обрабатывает одну строку:
// разбивает её по разделителю и выводит указанные колонки
// в соответствии с конфигурацией cfg.
func BuildMatcher(cfg *Config) func(string) {
	return func(s string) {
		if !strings.ContainsRune(s, cfg.Delim) {
			return
		}

		splitted := strings.Split(s, string(cfg.Delim))
		fields := make([]string, 0, len(cfg.ColIndexes))
		for _, colIdx := range cfg.ColIndexes {
			if colIdx >= len(splitted) {
				continue
			}
			fields = append(fields, splitted[colIdx])
		}
		if len(fields) == 0 {
			return
		}
		fmt.Println(strings.Join(fields, string(cfg.Delim)))
	}
}
