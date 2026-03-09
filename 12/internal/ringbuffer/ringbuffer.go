package ringbuffer

// Entry представляет одну строку входного потока,
// сохранённую вместе с её номером.
//
// Используется для хранения строк в кольцевом буфере
// и подавления повторного вывода.
type Entry struct {
	No   int
	Text string
}

// RingBuffer реализует кольцевой буфер фиксированного размера
// для хранения последних строк входного потока.
//
// Используется для реализации контекста до совпадения (-B).
type RingBuffer struct {
	buf        []Entry
	head, tail int
	full       bool
}

// Create создаёт кольцевой буфер заданной ёмкости.
// Если  размер буфера равен 0, возвращается nil.
func Create(cap int) *RingBuffer {
	if cap == 0 {
		return nil
	}
	return &RingBuffer{
		buf: make([]Entry, cap),
	}
}

// Push добавляет элемент в кольцевой буфер.
//
// При переполнении буфера самый старый элемент
// автоматически затирается.
func (b *RingBuffer) Push(e Entry) {
	if b == nil {
		return
	}

	b.buf[b.head] = e
	if b.full {
		b.tail = (b.tail + 1) % (len(b.buf))
	}
	b.head = (b.head + 1) % (len(b.buf))
	b.full = b.head == b.tail
}

// Snapshot возвращает срез элементов буфера
// в порядке от самого старого к самому новому.
//
// Состояние буфера при этом не изменяется.
func (b *RingBuffer) Snapshot() []Entry {
	if b == nil {
		return nil
	}
	n := b.len()
	if n == 0 {
		return nil
	}
	out := make([]Entry, n)
	for i := 0; i < n; i++ {
		out[i] = b.buf[(b.tail+i)%len(b.buf)]
	}
	return out
}

func (b *RingBuffer) len() int {
	if b == nil {
		return 0
	}
	if b.full {
		return len(b.buf)
	}
	if b.head >= b.tail {
		return b.head - b.tail
	}
	return len(b.buf) - b.tail + b.head
}
