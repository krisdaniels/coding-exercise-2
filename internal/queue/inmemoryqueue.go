package queue

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var _ Queue = (*inMemoryQueue)(nil)

type inMemoryQueue struct {
	commands []string
}

func NewInMemeoryQueue(read, write bool) Queue {
	commands := []string{}
	if write {
		commands = append(commands, "addItem('key%d','value%d')")
		commands = append(commands, "deleteItem('key%d')")
	}
	if read {
		commands = append(commands, "getItem('key%d')")
		commands = append(commands, "getAllItems()")
	}
	return &inMemoryQueue{
		commands: commands,
	}
}

func (q *inMemoryQueue) OpenConsumer() error {
	return nil
}

func (q *inMemoryQueue) Open() error {
	return nil
}

func (q *inMemoryQueue) Close() error {
	return nil
}

func (q *inMemoryQueue) Publish(_ string, _ time.Duration) error {
	return nil
}

func (q *inMemoryQueue) ReadNext() (string, error) {
	i := rand.Intn(32) % len(q.commands)
	j := rand.Intn(1024)

	cmd := q.commands[i]
	cnt := strings.Count(cmd, "%d")
	if cnt == 2 {
		cmd = fmt.Sprintf(cmd, j, j)
	} else if cnt == 1 {
		cmd = fmt.Sprintf(cmd, j)
	}

	// rate limit the in memory queue to max 50000 /s
	<-time.After(time.Microsecond * 20)
	return cmd, nil
}
