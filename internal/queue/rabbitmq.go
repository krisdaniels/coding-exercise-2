package queue

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// based on the hello world example on the rabbitmq website, kept default options to keep things simple

var _ Queue = (*rabbitMQ)(nil)

type rabbitMQ struct {
	endpoint  string
	queueName string

	closing bool
	conn    *amqp.Connection
	ch      *amqp.Channel
	queue   amqp.Queue
	msgs    <-chan amqp.Delivery
}

func NewRabbitMQ(endpoint, queueName string) Queue {
	return &rabbitMQ{
		endpoint:  endpoint,
		queueName: queueName,
	}
}

func (q *rabbitMQ) OpenConsumer() error {
	msgs, err := q.ch.Consume(
		q.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}
	q.msgs = msgs

	return err
}

func (q *rabbitMQ) Open() error {
	conn, err := amqp.Dial(q.endpoint)
	if err != nil {
		return err
	}
	q.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	q.ch = ch

	// use defaults from example for queue creation, can be extended to use config
	queue, err := ch.QueueDeclare(
		q.queueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}
	q.queue = queue

	return nil
}

func (q *rabbitMQ) Close() error {
	q.closing = true
	if q.ch != nil {
		q.ch.Close()
	}
	if q.conn != nil {
		q.conn.Close()
	}
	return nil
}

func (q *rabbitMQ) ReadNext() (string, error) {
	msg, ok := <-q.msgs
	if ok {
		return string(msg.Body), nil
	}
	if q.closing {
		return "", ErrorQueueClosing
	}

	return "", ErrorReadingMessage
}

func (q *rabbitMQ) Publish(command string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := q.ch.PublishWithContext(
		ctx,
		"",           // exchange
		q.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(command),
		})

	return err
}
