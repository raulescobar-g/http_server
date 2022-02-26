package main

import "fmt"

/*
Practice file for me to learn go basic singly linked list

*/
type Node struct {
	val  interface{}
	next *Node
}

type List struct {
	head *Node
	size int
}

func (list *List) insert(key interface{}) {
	list.size += 1
	n := Node{
		val:  key,
		next: nil,
	}

	if list.head == nil {
		list.head = &n
	} else {
		n.next = list.head
		list.head = &n
	}
}

func (list *List) print() {
	iter := list.head
	for i := 0; i < list.size; i++ {
		fmt.Printf("%d -> ", iter.val)
		iter = iter.next
	}
	fmt.Printf("\n")
}

func (list *List) remove(index int) int {
	if index >= list.size {
		return -1
	}

	if index == 0 {
		if list.head.next != nil {
			list.head = list.head.next
			list.size -= 1
			return 0
		}
		return -1
	}

	iter := list.head
	lag := iter

	for i := 0; i < index; i++ {
		lag = iter
		iter = iter.next
	}

	lag.next = iter.next
	list.size -= 1
	return 0
}
