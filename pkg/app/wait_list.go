package app

import "sync"

// CustomerInfo stores relevant information
// about a customer that needs to be placed on a wait list
type CustomerInfo struct {
	Msisdn string
	Ticket string
}

// WaitList is a queue data structure to store customer infos
// in a first-come first-served order
type WaitList struct {
	Items []*CustomerInfo
	lock  sync.RWMutex
}

// NewWaitList returns a pointer to a new waitlist
func NewWaitList() *WaitList {
	return &WaitList{
		Items: []*CustomerInfo{},
	}
}

// Enqueue puts a customer info at the end of the waiting list
func (wl *WaitList) Enqueue(c *CustomerInfo) {
	wl.lock.Lock()
	wl.Items = append(wl.Items, c)
	wl.lock.Unlock()
}

// Deque returns the customer info at the head of the waiting list
func (wl *WaitList) Deque() *CustomerInfo {
	wl.lock.Lock()
	customerInfo := wl.Items[0]
	wl.Items = wl.Items[1:len(wl.Items)]
	wl.lock.Unlock()

	return customerInfo
}

// IsEmpty returns true if the waiting list is empty
// or false otherwise
func (wl *WaitList) IsEmpty() bool {
	return len(wl.Items) == 0
}

// Size returns the number of customer info in the waiting list
func (wl *WaitList) Size() int {
	return len(wl.Items)
}
