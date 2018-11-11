package app

import (
	"fmt"
	"math/rand"
	"sync"

	"go.uber.org/zap"
)

// Rubix keeps track of internal state of the system in realtime
// 'WaitLists' is a mapping between queues defined in the db and
// their corresponding backing waitlist.
//
// 'nextTicketNumber' tracks the ticket number to be issued to the
// next customer who joins a queue
type Rubix struct {
	WaitLists        map[int]*WaitList
	nextTicketNumber int
	lock             sync.RWMutex
	logger           *zap.Logger
}

// NewRubix returns a pointer to a new State
func NewRubix(waitLists map[int]*WaitList, logger *zap.Logger) *Rubix {
	return &Rubix{
		WaitLists:        waitLists,
		nextTicketNumber: 1,
		logger:           logger,
	}
}

// Reset clears all application data
func (r *Rubix) Reset() {
	r.WaitLists = map[int]*WaitList{}
	r.lock.Lock()
	r.nextTicketNumber = 1
	r.lock.Unlock()
	r.logger.Info("application state reset", zap.Int("next_ticket_number", r.nextTicketNumber), zap.Any("wait_lists", r.WaitLists))
}

// GenerateTicket returns the next ticket identifier
func (r *Rubix) GenerateTicket() string {
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	ticket := fmt.Sprintf("%s%03d", letters[rand.Intn(len(letters))], r.nextTicketNumber)
	r.nextTicketNumber++

	return ticket
}

// AddCustomerToWaitList adds a customer info to the tail of a waitlist
// identied by the given queueId
func (r *Rubix) AddCustomerToWaitList(queueID int, msisdn, ticket string) {
	customerInfo := &CustomerInfo{Msisdn: msisdn, Ticket: ticket}

	_, ok := r.WaitLists[queueID]
	if !ok {
		r.logger.Info("creating waitlist for new queue", zap.Int("queue_id", queueID))
		r.WaitLists[queueID] = NewWaitList()
	}

	r.WaitLists[queueID].Enqueue(customerInfo)
	r.logger.Info("customer added to queue", zap.Any("customer_info", customerInfo), zap.Int("queueID", queueID))
}
