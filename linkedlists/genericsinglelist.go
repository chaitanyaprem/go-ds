package linkedlists

import (
	"errors"
	"fmt"
	"strconv"
)

/*
Wrote this in order to practice generics in golang
and understand the limitations of it.
*/

type GenericSingleList[T any] struct {
	head *GenericSinglelistNode[T]
	len  int
}

type GenericSinglelistNode[T any] struct {
	Value T
	Next  *GenericSinglelistNode[T]
}

func TraverseTill[T comparable](list *GenericSingleList[T], index int) *GenericSinglelistNode[T] {
	if index == -1 {
		index = list.len
	}
	tNode := list.head
	for i := 1; tNode != nil && tNode.Next != nil && i < index; tNode = tNode.Next {
		i++
	}
	return tNode
}

/*
Add adds a value to the end of the list and returns the index at which it was added.
Index starts at 1.
*/
func Add[T comparable](list *GenericSingleList[T], value T) {
	var node GenericSinglelistNode[T]
	node.Value = value
	if list.head == nil {
		list.head = &node
	} else {
		lastNode := TraverseTill(list, -1)
		lastNode.Next = &node
	}
	list.len++
}

/*
This method deletes the first entry with the value provided.
Returns error in case it can't find an entry with value passed.
*/
func Delete[T comparable](list *GenericSingleList[T], value T) error {
	tNode := list.head
	i := 1
	var prev *GenericSinglelistNode[T]
	for ; i <= list.len; i++ {
		if tNode.Value == value {
			break
		}
		if tNode.Next == nil {
			return errors.New("could not find entry with value in the list")
		}
		prev = tNode
		tNode = tNode.Next
	}
	if i == 1 {
		list.head = list.head.Next
		list.len--
		return nil
	} else if i == list.len {
		prev.Next = nil
	} else {
		prev.Next = tNode.Next
	}
	fmt.Printf("Entry being deleted at index %d is %+v \n", i, tNode)
	list.len--
	return nil
}

/*
This method deletes an entry at index provided.
Returns the value at the index if delete is successful, else error in case of invalid index.
*/
func DelEntryAtIndex[T comparable](list *GenericSingleList[T], index int) (T, error) {
	var defValue T
	if index == -1 || index == 0 || index > list.len {
		return defValue, errors.New("invalid index " + strconv.Itoa(index) + " passed")
	}
	tNode := list.head
	if index == 1 {
		list.head = list.head.Next
	} else {
		var prev *GenericSinglelistNode[T]
		for i := 1; i < index; i++ {
			if tNode.Next == nil {
				break
			}
			prev = tNode
			tNode = tNode.Next
		}
		prev.Next = tNode.Next
	}
	fmt.Printf("Entry at index %d is %+v \n", index, tNode)
	list.len--
	return tNode.Value, nil
}

/*
This method prints the contents of the list in order
*/
func PrintList[T comparable](list *GenericSingleList[T]) {
	fmt.Println("Printing all elements of list of length :", list.len)
	tNode := list.head
	for i := 0; tNode != nil; tNode = tNode.Next {
		fmt.Println("Index: ", i, ", Value: ", tNode.Value)
		i++
	}
}

func Length[T comparable](list *GenericSingleList[T]) int {
	return list.len
}
