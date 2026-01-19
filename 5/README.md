# Задача l2.5

## Задание
Что выведет программа?

Объясните вывод программы.
```go
package main

type customError struct {
  msg string
}

func (e *customError) Error() string {
  return e.msg
}

func test() *customError {
  // ... do something
  return nil
}

func main() {
  var err error
  err = test()
  if err != nil {
    println("error")
    return
  }
  println("ok")
}
```

## Решение

Программа выведет:
```
error
```
### Объяснение

Функция `test` возвращает значение типа `*customError` и при выходе из функции `test` в `main` переменной интерфейсного типа `error` присваивается это значение и `error` уже не будет равен `nil`, как мы знаем из задачи [l2.2](https://github.com/vvigg0/wbtech-l2/tree/main/2).

В текущей реализации функции `test` всегда будет выводиться `error`

#### Как исправить

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() error {
	// ... do something
	// if ... {
	// 	return &customError{}
	// }
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

```
Функция `test` теперь возвращает интерфейс `error`, когда произойдет какая-то ошибка функция явно вернет `error` с типом `*customError` и будет `err!=nil`, если же ошибки в `test` нет, возвращается `nil` интерфейс и программа выведет `ok`