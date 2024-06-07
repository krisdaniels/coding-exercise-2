package queue

import (
	"errors"
	"time"
)

var (
	ErrorQueueClosing   = errors.New("closing")
	ErrorReadingMessage = errors.New("error reading message")
)

type Queue interface {
	Open() error
	OpenConsumer() error
	Close() error
	ReadNext() (string, error)
	Publish(command string, timeout time.Duration) error
}
