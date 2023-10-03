package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{
		v,
		l.first,
		nil,
	}

	if l.first == nil {
		l.first = i
		l.last = i
		l.len++

		return i
	}

	l.first.Prev = i
	l.first = i
	l.len++

	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.last == nil {
		return l.PushFront(v)
	}

	i := &ListItem{
		v,
		nil,
		l.last,
	}

	l.last.Next = i
	l.last = i
	l.len++

	return i
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.first == i && l.last == i:
		l.first = nil
		l.last = nil
		l.len--

		return
	case l.first == i:
		i.Next.Prev = nil
		l.first = i.Next
		l.len--

		return
	case l.last == i:
		i.Prev.Next = nil
		l.last = i.Prev
		l.len--

		return
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		l.len--

		return
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.first == i {
		return
	}

	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Next == nil {
		l.last = i.Prev
	}
	i.Prev = nil
	i.Next = l.first
	l.first.Prev = i
	l.first = i
}

func NewList() List {
	return new(list)
}
