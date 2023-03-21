package linkedlists

import (
	"errors"
	"fmt"
	"strconv"
)

type SingleList struct {
	head *singleListNode
	len  int
}

type singleListNode struct {
	Value int
	Next  *singleListNode
}

/*
This method returns the node present at the index in the list.
If index is -1 , end of list is returned.
If index is greater than lenght of list, end of list is returned.
*/
func (list *SingleList) TraverseTill(index int) *singleListNode {
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
func (list *SingleList) Add(value int) {
	var node singleListNode
	node.Value = value
	if list.head == nil {
		list.head = &node
	} else {
		lastNode := list.TraverseTill(-1)
		lastNode.Next = &node
	}
	list.len++
}

/*
This method deletes the first entry with the value provided.
Returns error in case it can't find an entry with value passed.
*/
func (list *SingleList) Delete(value int) error {
	tNode := list.head
	i := 1
	var prev *singleListNode
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
func (list *SingleList) DelEntryAtIndex(index int) (int, error) {
	if index <= 0 || index > list.len {
		return -1, errors.New("invalid index " + strconv.Itoa(index) + " passed")
	}
	tNode := list.head
	if index == 1 {
		list.head = list.head.Next
	} else {
		var prev *singleListNode
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
func (list *SingleList) PrintList() {
	fmt.Println("Printing all elements of list of length :", list.len)
	tNode := list.head
	for i := 0; tNode != nil; tNode = tNode.Next {
		fmt.Println("Index: ", i, ", Value: ", tNode.Value)
		i++
	}
}

func (list *SingleList) Length() int {
	return list.len
}

func (list *SingleList) GetValueAt(index int) (int, error) {
	if index <= 0 || index > list.len {
		return -1, errors.New("invalid index " + strconv.Itoa(index) + " passed")
	}
	tNode := list.TraverseTill(index)
	return tNode.Value, nil
}
