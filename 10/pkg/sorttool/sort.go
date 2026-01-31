package sorttool

import (
	"sort"
)

// CheckSort проверяет, отсортирован ли срез строк в соответствии с функцией сравнения cmp.
// Если unique равно true, одинаковые соседние элементы считаются ошибкой.
// В случае ошибки функция возвращает false и индекс первого некорректного элемента
func CheckSort(s []string, cmp Comparator, unique bool) (bool, int) {
	for i := 1; i < len(s); i++ {
		c := cmp(s[i-1], s[i])
		if c > 0 {
			return false, i
		}
		if unique && c == 0 {
			return false, i
		}
	}
	return true, -1
}

// Sort сортирует срез строк на месте, используя стабильную сортировку
// и предоставленный компаратор.
func Sort(s []string, cmp Comparator) {
	sort.SliceStable(s, func(i int, j int) bool {
		return cmp(s[i], s[j]) < 0
	})
}

// UniqueSorted удаляет смежные повторяющиеся элементы из отсортированного среза.
// Входной срез должен быть предварительно отсортирован с использованием того же компаратора.
func UniqueSorted(s []string, cmp Comparator) []string {
	if len(s) == 0 {
		return s
	}
	out := make([]string, 0, len(s))
	out = append(out, s[0])
	for i := 1; i < len(s); i++ {
		if cmp(s[i-1], s[i]) != 0 {
			out = append(out, s[i])
		}
	}
	return out
}
