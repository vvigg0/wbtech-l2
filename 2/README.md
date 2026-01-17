# Задача l2.2

## Задание
Что выведет программа?

Объяснить порядок выполнения defer функций и итоговый вывод.
```go
package main

import "fmt"

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}

func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```
## Решение

Программа выведет `2 1`. 

### Объяснение
`defer` выполняется после вычисления `return-value`, но до фактического выхода из функции

В функции `test` используется именованный возврат: при достижении `return`, `x` не сохраняется как `return-value` и `defer` его изменяет. На выходе `x=2`

В функции `anothertest` возвращаемое значение - неименованное: при достижении `return`, `x` сохраняется как `return-value`, после чего `defer` изменяет локальную переменную `x`, что уже не влияет на `return-value`