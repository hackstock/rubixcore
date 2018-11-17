package app

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// NandiSMSWorker sends text messages through Nandi Mobile
type NandiSMSWorker struct {
	brokerConn    *amqp.Connection
	gatewayConfig *SmsGatewayConfig
	logger        *zap.Logger
}

// NewNandiSMSWorker returns a NandiSMSWorker
func NewNandiSMSWorker(brokerConn *amqp.Connection, config *SmsGatewayConfig, logger *zap.Logger) NandiSMSWorker {
	return NandiSMSWorker{
		brokerConn:    brokerConn,
		gatewayConfig: config,
		logger:        logger,
	}
}

// Run starts a goroutine that consumes messages on the
// sms_task_queue and makes HTTP calls to send SMS
// via Nandi Mobile's SMS gateway
func (worker NandiSMSWorker) Run(queueName string) {
	channel, err := worker.brokerConn.Channel()
	if err != nil {
		worker.logger.Info("failed creating channel for worker", zap.Error(err))
	}

	messages, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		worker.logger.Info("failed consuming messages on queue", zap.Error(err), zap.String("queue_name", queueName))
	}

	go func() {
		for data := range messages {
			payload := string(data.Body)
			worker.logger.Info("payload received", zap.String("payload", payload))
			details := strings.Split(payload, "#")
			msisdn := details[0]
			message := details[1]
			err := sendSMS(message, msisdn, worker.gatewayConfig)
			if err != nil {
				worker.logger.Error("failed sending SMS via Nandi", zap.Error(err))
			}
			data.Ack(true)
		}
	}()
}

func sendSMS(msg, to string, config *SmsGatewayConfig) error {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	form := url.Values{}
	form.Add("username", config.Username)
	form.Add("password", config.Password)
	form.Add("numbers", to)
	form.Add("message", msg)
	form.Add("from", config.SenderID)

	url := fmt.Sprintf("https://infoline.nandiclient.com/%s/campaigns/sendmsg", config.Username)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
