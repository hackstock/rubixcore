package app

import "github.com/streadway/amqp"

// SMSPublisher publishes SMS payloads onto
// task queues to be consumed by SMS workers
type SMSPublisher struct {
	brokerConn *amqp.Connection
}

// NewSMSPublisher returns a pointer to a new SMSPublisher
func NewSMSPublisher(conn *amqp.Connection) SMSPublisher {
	return SMSPublisher{
		brokerConn: conn,
	}
}

// Publish publishes sms onto the given queueName on
// the message broker connection
func (p SMSPublisher) Publish(sms, queueName string) error {
	channel, err := p.brokerConn.Channel()
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(
		smsTaskQueue, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(sms),
		})

	return err
}
