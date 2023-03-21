package linkedlists

import (
	"testing"
)

func TestListAdd(t *testing.T) {

	values := []int{1, 2, 3, 10, 20, 100, 55, 130}
	var sList SingleList
	for _, v := range values {
		sList.Add(v)
	}
	//Verify values at indices

	for i, v := range values {
		listIndex := i + 1
		listVal, err := sList.GetValueAt(listIndex)
		if err != nil {
			t.Fatalf("Add Failure. Invalid index passed while fetching entry")
		}
		if listVal != v {
			sList.PrintList()
			t.Fatalf("Add Failure. Value %d in the list at index %d is different than inserted value of %d", listVal, i, v)
		}
	}
	//sList.PrintList()
}

func TestListDelete(t *testing.T) {

	values := []int{1, 2, 3, 10, 20, 100, 55, 130}
	var sList SingleList
	for _, v := range values {
		sList.Add(v)
	}

	for _, v := range values {
		err := sList.Delete(v)
		if err != nil {
			sList.PrintList()
			t.Fatalf("Delete failed from list for value %d", v)
		}
	}

}

func TestListDeleteAtIndex(t *testing.T) {
	values := []int{1, 2, 3, 10, 20, 100, 55, 130}
	var sList SingleList
	for _, v := range values {
		sList.Add(v)
	}
	for i, v := range values {
		listVal, err := sList.DelEntryAtIndex(1)
		if err != nil {
			sList.PrintList()
			t.Fatalf("Delete failed with error %s using DelEntryAtIndex (%d)", err.Error(), i)
		}
		if listVal != v {
			sList.PrintList()
			t.Fatalf("Delete failed while using DelEntryAtIndex (%d)", i)
		}
	}

}
