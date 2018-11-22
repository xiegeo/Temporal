package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database"
	"github.com/streadway/amqp"
)

func (qm *Manager) setupLogging() error {
	logFileName := fmt.Sprintf("/var/log/temporal/%s_serice.log", qm.QueueName)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	logger := log.New()
	logger.Out = logFile
	qm.Logger = logger
	qm.Logger.Info("Logging initialized")
	return nil
}

func (qm *Manager) parseQueueName(queueName string) error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}
	qm.QueueName = fmt.Sprintf("%s+%s", host, queueName)
	return nil
}

// Initialize is used to connect to the given queue, for publishing or consuming purposes
func Initialize(queueName, connectionURL string, publish, service bool) (*Manager, error) {
	conn, err := setupConnection(connectionURL)
	if err != nil {
		return nil, err
	}
	qm := Manager{Connection: conn}
	if err := qm.OpenChannel(); err != nil {
		return nil, err
	}

	qm.QueueName = queueName
	qm.Service = queueName

	if service {
		err = qm.setupLogging()
		if err != nil {
			return nil, err
		}
	}
	// Declare Non Default exchanges for the particular queue
	switch queueName {
	case IpfsPinQueue:
		err = qm.parseQueueName(queueName)
		if err != nil {
			return nil, err
		}
		err = qm.DeclareIPFSPinExchange()
		if err != nil {
			return nil, err
		}
		qm.ExchangeName = PinExchange
		if publish {
			return &qm, nil
		}
	case IpfsKeyCreationQueue:
		err = qm.parseQueueName(queueName)
		if err != nil {
			return nil, err
		}
		err = qm.DeclareIPFSKeyExchange()
		if err != nil {
			return nil, err
		}
		qm.ExchangeName = IpfsKeyExchange
		if publish {
			return &qm, nil
		}
	}
	if err := qm.DeclareQueue(); err != nil {
		return nil, err
	}
	return &qm, nil
}

func setupConnection(connectionURL string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// OpenChannel is used to open a channel to the rabbitmq server
func (qm *Manager) OpenChannel() error {
	ch, err := qm.Connection.Channel()
	if err != nil {
		return err
	}
	if qm.Logger != nil {
		qm.LogInfo("channel opened")
	}
	qm.Channel = ch
	return nil
}

// DeclareQueue is used to declare a queue for which messages will be sent to
func (qm *Manager) DeclareQueue() error {
	// we declare the queue as durable so that even if rabbitmq server stops
	// our messages won't be lost
	q, err := qm.Channel.QueueDeclare(
		qm.QueueName, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}
	if qm.Logger != nil {
		qm.LogInfo("queue declared")
	}
	qm.Queue = &q
	return nil
}

// ConsumeMessage is used to consume messages that are sent to the queue
// Question, do we really want to ack messages that fail to be processed?
// Perhaps the error was temporary, and we allow it to be retried?
func (qm *Manager) ConsumeMessage(consumer, dbPass, dbURL, dbUser string, cfg *config.TemporalConfig) error {
	db, err := database.OpenDBConnection(database.DBOptions{
		User:     cfg.Database.Username,
		Password: cfg.Database.Password,
		Address:  cfg.Database.URL,
		Port:     cfg.Database.Port,
	})
	if err != nil {
		return err
	}

	// ifs the queue is using an exchange, we will need to bind the queue to the exchange
	switch qm.ExchangeName {
	case PinRemovalExchange, PinExchange, IpfsKeyExchange:
		if err = qm.Channel.QueueBind(
			qm.QueueName,    // name of the queue
			"",              // routing key
			qm.ExchangeName, // exchange
			false,           // noWait
			nil,             // arguments
		); err != nil {
			return err
		}
	default:
		break
	}

	// consider moving to true for auto-ack
	msgs, err := qm.Channel.Consume(
		qm.QueueName, // queue
		consumer,     // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	// check the queue name
	switch qm.Service {
	// only parse database file requests
	case DatabaseFileAddQueue:
		return qm.ProcessDatabaseFileAdds(msgs, db)
	case IpfsPinQueue:
		return qm.ProccessIPFSPins(msgs, db, cfg)
	case IpfsFileQueue:
		return qm.ProccessIPFSFiles(msgs, cfg, db)
	case EmailSendQueue:
		return qm.ProcessMailSends(msgs, cfg)
	case IpnsEntryQueue:
		return qm.ProcessIPNSEntryCreationRequests(msgs, db, cfg)
	case IpfsKeyCreationQueue:
		return qm.ProcessIPFSKeyCreation(msgs, db, cfg)
	case IpfsClusterPinQueue:
		return qm.ProcessIPFSClusterPins(msgs, cfg, db)
	default:
		return errors.New("invalid queue name")
	}
}

//PublishMessageWithExchange is used to publish a message to a given exchange
func (qm *Manager) PublishMessageWithExchange(body interface{}, exchangeName string) error {
	switch exchangeName {
	case PinExchange:
		break
	case PinRemovalExchange:
		break
	case IpfsKeyExchange:
		break
	default:
		return errors.New("invalid exchange name provided")
	}
	bodyMarshaled, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if err = qm.Channel.Publish(
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         bodyMarshaled,
		},
	); err != nil {
		return err
	}
	return nil
}

// PublishMessage is used to produce messages that are sent to the queue, with a worker queue (one consumer)
func (qm *Manager) PublishMessage(body interface{}) error {
	// we use a persistent delivery mode to combine with the durable queue
	bodyMarshaled, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if err = qm.Channel.Publish(
		"",            // exchange
		qm.Queue.Name, // routing key
		false,         // mandatory
		false,         //immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         bodyMarshaled,
		},
	); err != nil {
		return err
	}
	return nil
}

// Close is used to close our queue connection
func (qm *Manager) Close() error {
	return qm.Connection.Close()
}
