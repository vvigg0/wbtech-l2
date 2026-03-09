package cutanalog

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

var errGreaterZero error = errors.New("Номер колонки должен быть больше нуля")
var errNotNumber error = errors.New("Во флаге f допустимы только числа")

// Config содержит параметры работы утилиты cut:
// индексы колонок, которые нужно вывести;
// символ-разделитель между колонками;
// флаг -s: нужно ли игнорировать строки без разделителя
type Config struct {
	ColIndexes []int
	Delim      rune
	Sep        bool
}

// ParseFields парсит значение флага -f, валидирует его
// и заполняет Config.ColIndexes индексами колонок (с нуля).
// Поддерживает одиночные значения и диапазоны (например: 1,3-5).
func ParseFields(f string, cfg *Config) error {
	if f == "" {
		return errors.New("Номер колонки обязателен")
	}
	fields := strings.Split(f, ",")
	for _, str := range fields {
		if strings.Contains(str, "-") {
			parts := strings.Split(str, "-")
			if len(parts) > 2 {
				return errors.New("Диапазон должен состоять из двух чисел (2-7)")
			}
			l := parts[0]
			r := parts[1]
			lIdx, err := fieldToColIdx(l)
			if err != nil {
				return err
			}
			rIdx, err := fieldToColIdx(r)
			if err != nil {
				return err
			}
			if lIdx > rIdx {
				return errors.New("Левая граница должна быть меньше или равна правой границе")
			}
			for lIdx <= rIdx {
				cfg.ColIndexes = append(cfg.ColIndexes, lIdx)
				lIdx++
			}
		} else {
			colIdx, err := fieldToColIdx(str)
			if err != nil {
				return err
			}
			cfg.ColIndexes = append(cfg.ColIndexes, colIdx)
		}
	}
	sort.Ints(cfg.ColIndexes)
	cfg.ColIndexes = dedupSortedInts(cfg.ColIndexes)
	return nil
}

func fieldToColIdx(s string) (int, error) {
	intS, err := strconv.Atoi(s)
	if err != nil {
		return -1, errNotNumber
	}
	if intS < 1 {
		return -1, errGreaterZero
	}
	return intS - 1, nil
}
func dedupSortedInts(s []int) []int {
	if len(s) == 0 {
		return s
	}
	j := 0
	for i := 1; i < len(s); i++ {
		if s[i] != s[j] {
			j++
			s[j] = s[i]
		}
	}
	return s[:j+1]

}
