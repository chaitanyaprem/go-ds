package main

import (
	"fmt"
	"go-ds/ringbuffer"
)

func main() {
	/* 	fmt.Println("Testing Linked List")
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

	   	sList.PrintList() */
	var buffer ringbuffer.RingBuffer
	buffer.Initialize(10)
	fmt.Println("Initialized ringBuffer with size 10")
	data := []byte{102, 97}
	for i := 0; i < 5; i++ {
		err := buffer.Write(data)
		if err != nil {
			fmt.Println("Failed to Write to buffer due to error: ", err)
			break
		}
		fmt.Println("Wrote ", len(data), " bytes successfully to ringBuffer.")
	}
	err := buffer.Write(data)
	if err != nil {
		fmt.Println("Failed to Write to buffer due to error: ", err)
	}
	buffer.Print()
	readBytes, err := buffer.Read(6)
	if err != nil {
		fmt.Println("Error reading from buffer ", err)
	}
	fmt.Printf("Read %d bytes from buffer. Bytes: %+v \n", len(readBytes), readBytes)
	buffer.Print()
	err = buffer.Write(data)
	if err != nil {
		fmt.Println("Failed to Write to buffer due to error: ", err)
	}
	fmt.Println("Wrote ", len(data), " bytes successfully")
	buffer.Print()
	readBytes, err = buffer.Read(8)
	if err != nil {
		fmt.Println("Error reading from buffer ", err)
	}
	fmt.Printf("Read %d bytes from buffer. Bytes: %+v \n", len(readBytes), readBytes)
	buffer.Print()
	for i := 0; i < 6; i++ {
		err = buffer.Write(data)
		if err != nil {
			fmt.Println("Failed to Write to buffer due to error: ", err)
			break
		}
	}
	buffer.Print()
}
