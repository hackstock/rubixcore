package app

import (
	"reflect"
	"testing"
)

func TestWaitList(t *testing.T) {
	waitList := NewWaitList()
	if !waitList.IsEmpty() {
		t.Errorf("expected waitlist to be empty, got non empty")
	}

	if waitList.Size() != 0 {
		t.Errorf("expected a size of 0, got %d", waitList.Size())
	}

	ciOne := CustomerInfo{ID: 1, Msisdn: "+233200662782", Ticket: "A100"}
	ciTwo := CustomerInfo{ID: 2, Msisdn: "+233200662783", Ticket: "A200"}
	ciThree := CustomerInfo{ID: 3, Msisdn: "+233200662789", Ticket: "A300"}

	waitList.Enqueue(&ciOne)
	if waitList.Size() != 1 {
		t.Errorf("expected a size of 1, got %d", waitList.Size())
	}

	waitList.Enqueue(&ciTwo)
	if waitList.Size() != 2 {
		t.Errorf("expected a size of 1, got %d", waitList.Size())
	}

	waitList.Enqueue(&ciThree)
	if waitList.Size() != 3 {
		t.Errorf("expected a size of 1, got %d", waitList.Size())
	}

	got := waitList.Deque()
	if !reflect.DeepEqual(got, &ciOne) {
		t.Errorf("expected %v, got %v", ciOne, got)
	}
}
