package main

import (
	"fmt"
	"go-ds/linkedlists"
)

func main() {
	fmt.Println("Testing Linked List")
	var sList linkedlists.SingleList
	sList.Add(1)
	sList.Add(2)
	sList.Add(3)
	fmt.Println("Length of list is ", sList.Length())
	delAtIndex := 2
	val, err := sList.DelEntryAtIndex(delAtIndex)
	if err != nil {
		fmt.Println("Error deleting entry ", err)
		err = nil
	}
	fmt.Println("Entry deleted at index ", delAtIndex, " is ", val)
	fmt.Println("Length of list is ", sList.Length())
	sList.PrintList()

	err = sList.Delete(2)
	if err != nil {
		fmt.Println("Error deleting entry with value 2", err)
		err = nil
	}
	sList.PrintList()

	err = sList.Delete(3)
	if err != nil {
		fmt.Println("Error deleting entry with value 3", err)
		err = nil
	}
	sList.Delete(1)
	fmt.Println("Length of list is ", sList.Length())

	sList.PrintList()
}
