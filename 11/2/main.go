// В данном решении для определения анаграмм используется
// сигнатура слова, полученная путём сортировки его символов.
package main

import (
	"fmt"
	"sort"
	"strings"
)

func getSet(words []string) map[string][]string {
	m := make(map[string][]string)
	first := make(map[string]string)
	seen := make(map[string]struct{})
	for _, word := range words {
		lw := strings.ToLower(word)
		if _, ok := seen[lw]; ok {
			continue
		}
		seen[lw] = struct{}{}

		runes := []rune(lw)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		sig := string(runes)

		if _, ok := first[sig]; !ok {
			first[sig] = lw
		}
		m[sig] = append(m[sig], lw)
	}
	res := make(map[string][]string)
	for sig, v := range m {
		if len(v) < 2 {
			continue
		}
		sort.Strings(v)
		res[first[sig]] = v
	}
	return res
}
func main() {
	words := []string{"пятка", "пятак", "тяпка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(getSet(words))
}
