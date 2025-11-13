package redict

type node struct {
	v    []byte
	next *node
	prev *node
}

type list struct {
	head   *node
	tail   *node
	length uint32
}

func newList() *list {
	return &list{}
}

func (l *list) insertHead(v []byte) {
	n := &node{
		v:    v,
		next: l.head,
		prev: nil,
	}

	if l.head == nil {
		l.tail = n
	} else {
		l.head.prev = n
	}

	l.head = n
	l.length++
}

func (l *list) insertTail(v []byte) {
	n := &node{
		v:    v,
		next: nil,
		prev: l.tail,
	}

	if l.tail == nil {
		l.head = n
	} else {
		l.tail.next = n
	}

	l.tail = n
	l.length++
}

func (l *list) popHead() []byte {
	if l.head == nil {
		return nil
	}

	n := l.head
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	} else {
		l.head.prev = nil
		n.next = nil
	}

	l.length--
	return n.v
}

func (l *list) popTail() []byte {
	if l.tail == nil {
		return nil
	}

	n := l.tail
	l.tail = l.tail.prev
	if l.tail == nil {
		l.head = nil
	} else {
		l.tail.next = nil
		n.prev = nil
	}

	l.length--
	return n.v
}

func (l *list) get(start, end int64) [][]byte {
	if start < 0 {
		start = int64(l.length) + start
	}

	end = min(end, int64(l.length-1))
	if end < 0 {
		end = int64(l.length) + end
	}

	if start > end {
		return nil
	}

	n := l.head
	for range start {
		n = n.next
	}
	length := end - start + 1
	b := make([][]byte, length)
	for i := range length {
		if n == nil {
			break
		}

		b[i] = n.v
		n = n.next
	}

	return b
}

func (l *list) trim(start, end int64) {
	if l.head == nil {
		return
	}

	if start < 0 {
		start = int64(l.length) + start
	}

	if end < 0 {
		end = int64(l.length) + end
	}

	defer func() { l.length = uint32(max(0, int64(l.length)-start-(int64(l.length-1)-end))) }()

	for range start {
		l.head = l.head.next
		if l.head == nil {
			l.tail = nil
			return
		}
		l.head.prev = nil
	}

	for range max(0, int64(l.length-1)-end) {
		l.tail = l.tail.prev
		if l.tail == nil {
			l.head = nil
			return
		}
		l.tail.next = nil
	}
}
