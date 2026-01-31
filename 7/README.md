# Задача l2.7

## Задание
Что выведет программа?

Объяснить работу конвейера с использованием select.
```go
package main

import (
  "fmt"
  "math/rand"
  "time"
)

func asChan(vs ...int) <-chan int {
  c := make(chan int)
  go func() {
    for _, v := range vs {
      c <- v
      time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
    }
  close(c)
}()
  return c
}

func merge(a, b <-chan int) <-chan int {
  c := make(chan int)
  go func() {
    for {
      select {
        case v, ok := <-a:
          if ok {
            c <- v
          } else {
            a = nil
          }
        case v, ok := <-b:
          if ok {
            c <- v
          } else {
            b = nil
          }
        }
        if a == nil && b == nil {
          close(c)
          return
        }
     }
   }()
  return c
}

  func main() {
    rand.Seed(time.Now().Unix())
    a := asChan(1, 3, 5, 7)
    b := asChan(2, 4, 6, 8)
    c := merge(a, b)
    for v := range c {
    fmt.Print(v)
  }
}
```

## Решение

### Вывод

Программа выведет числа от 1 до 8 в случайном порядке. Например: 12345678 или 12346578

### Объяснение
Принцип работы select в данном случае:

- `select` блокируется до тех пор, пока хотя бы один из каналов не будет готов
- При получении значения из какого-либо канала, это значение пишется в результирующий канал `c`
	- (Если готовы несколько каналов, выбирается случайный)
- Когда канал закрывается, `ok` становится `false`, и канал устанавливается `nil`
- `nil` Канал в `select` никогда не готов - это исключает его из выбора