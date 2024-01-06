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
	Key   Key
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	item.Next = l.head
	item.Prev = nil
	if l.head != nil {
		l.head.Prev = item
	}
	l.head = item
	if l.tail == nil {
		l.tail = item
	}
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	item.Next = nil
	item.Prev = l.tail
	if l.tail != nil {
		l.tail.Next = item
	}
	l.tail = item
	if l.head == nil {
		l.head = item
	}
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.head == i {
		l.head = i.Next
	}
	if l.tail == i {
		l.tail = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}
	if l.tail == i {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	i.Next = l.head
	l.head.Prev = i
	l.head = i
	l.head.Prev = nil
	l.tail.Next = nil
}

func NewList() List {
	return new(list)
}
