# Задача l2.3

## Задание
Что выведет программа?

Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.
```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

## Решение

Программа выведет:
```
nil
false
```

### Объяснение

В Go есть 2 структуры описывающие интерфейсы: `iface` и `eface`.

В данном примере интерфейсу `error` присваивается значение `nil` типа `*os.PathError`.

Инструкция `fmt.Println(err)` выводит `nil` (значение), но `fmt.Println(err == nil)` выводит `false`, потому что интерфейс равен `nil` только в том случае, если у него нет динамического типа (`iface.tab == nil`) и значения (`iface.data == nil`). Здесь `tab != nil`,а `data == nil`
#### `iface` (не пустой интерфейс)
```go
type iface struct{
  tab *itab
  data unsafe.Pointer // конкретное значение
}

type itab struct {
    inter *interfacetype // конкретный интерфейс
    _type *_type         // конкретный тип
    hash  uint32
    fun   [N]uintptr     // методы конкретного типа
}
```


#### `eface` (пустой интерфейс)
```go
type eface struct{
  _type *_type
  data  unsafe.Pointer
}
```

Отличие `iface` от `eface` заключается в том, что `iface` используется, когда есть методы, а `eface` когда методов нет. `iface` может хранить в себе только те значения, типы которых удовлетворяют самому интерфейсу, а удовлетворяют они когда имеют методы описанные в интерфейсе. `eface` может хранить любые значения: пустому интерфейсу удовлетворяет любой тип,потому что интерфейс не содержит методов.