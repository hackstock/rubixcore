package app

import (
	"fmt"
	"math/rand"
	"sync"

	"go.uber.org/zap"
)

const (
	smsTaskQueue = "sms_task_queue"
)

// Rubix keeps track of internal state of the system in realtime
// 'WaitLists' is a mapping between queues defined in the db and
// their corresponding backing waitlist.
//
// 'nextTicketNumber' tracks the ticket number to be issued to the
// next customer who joins a queue
type Rubix struct {
	waitLists        map[int64]*WaitList
	nextTicketNumber int
	lock             sync.RWMutex
	publisher        Publisher
	logger           *zap.Logger
}

// Publisher publishes messages to a queue in a message broker
type Publisher interface {
	Publish(sms, queueName string) error
}

// SMSWorker consumes messages from a queue in a message broker,
// parses the messages, and makes HTTP calls to send SMS via
// configured SMS gateways
type SMSWorker interface {
	Run(queueName string)
}

// NewRubix returns a pointer to a new State
func NewRubix(waitLists map[int64]*WaitList, publisher Publisher, logger *zap.Logger) *Rubix {
	return &Rubix{
		waitLists:        waitLists,
		nextTicketNumber: 1,
		publisher:        publisher,
		logger:           logger,
	}
}

// Reset clears all application data
func (r *Rubix) Reset() {
	r.waitLists = map[int64]*WaitList{}
	r.lock.Lock()
	r.nextTicketNumber = 1
	r.lock.Unlock()
	r.logger.Info("application state reset", zap.Int("next_ticket_number", r.nextTicketNumber), zap.Any("wait_lists", r.waitLists))
}

// RegisterSMSWorker starts a worker tasks that sends SMS to customers
func (r *Rubix) RegisterSMSWorker(worker SMSWorker) {
	worker.Run(smsTaskQueue)
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
func (r *Rubix) AddCustomerToWaitList(queueID int64, msisdn, ticket string) error {
	customerInfo := &CustomerInfo{Msisdn: msisdn, Ticket: ticket}

	_, ok := r.waitLists[queueID]
	if !ok {
		r.logger.Info("creating waitlist for new queue", zap.Int64("queue_id", queueID))
		r.waitLists[queueID] = NewWaitList()
	}

	msg := fmt.Sprintf("Ticket number %s. Kindly wait for your turn.", ticket)
	smsPayload := fmt.Sprintf("%s#%s", msisdn, msg)
	err := r.publisher.Publish(smsPayload, smsTaskQueue)
	if err != nil {
		return err
	}

	r.waitLists[queueID].Enqueue(customerInfo)
	r.logger.Info("customer added to queue", zap.Any("customer_info", customerInfo), zap.Int64("queueID", queueID))

	return nil
}

// NotifyNextCustomer deques a customer and notifies him of his
// turn to be served at a specific counter and returns the
// ID of the customer
func (r *Rubix) NotifyNextCustomer(queueID, counterID int64) int {
	customer := r.waitLists[queueID].Deque()
	r.logger.Info("customer notified of turn", zap.Any("customer", customer), zap.Int64("counter", counterID))
	return 1
}
