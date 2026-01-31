// В данном решении для группировки анаграмм используется
// частотная сигнатура слова — массив из 33 целых чисел,
// соответствующий количеству вхождений каждой буквы
// русского алфавита (включая «ё»).
package main

import (
	"fmt"
	"sort"
	"strings"
)

func getSet(words []string) map[string][]string {
	seen := make(map[string]struct{})
	m := make(map[[33]int][]string)
	for _, word := range words {
		word = strings.ToLower(word)
		if _, ok := seen[word]; ok {
			continue
		}
		seen[word] = struct{}{}
		freq := getFrequency(word)
		m[freq] = append(m[freq], word)
	}
	result := make(map[string][]string)
	for _, v := range m {
		if len(v) == 1 {
			continue
		}
		result[v[0]] = v
		sort.Strings(result[v[0]])
	}
	return result
}
func getFrequency(word string) [33]int {
	var freq [33]int
	for _, letter := range word {
		if letter == 'ё' {
			freq[32]++
			continue
		}
		freq[letter-'а']++
	}
	return freq
}
func main() {
	words := []string{"пятак", "пятка", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(getSet(words))
}
