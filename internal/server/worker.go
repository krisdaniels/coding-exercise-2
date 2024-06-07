package server

import (
	"coding-exercise/internal/commandparser"
	"coding-exercise/internal/orderedmap"
	"coding-exercise/internal/queue"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type Worker interface {
	Start()
	Stop()
}

var _ Worker = (*worker)(nil)

type worker struct {
	data         orderedmap.OrderedMap
	queue        queue.Queue
	parser       commandparser.CommandParser
	outputFolder string
	stop         chan bool
}

func newWorker(
	data orderedmap.OrderedMap,
	queue queue.Queue,
	parser commandparser.CommandParser,
	outputFolder string,
) Worker {
	return &worker{
		data:         data,
		queue:        queue,
		parser:       parser,
		outputFolder: outputFolder,
		stop:         make(chan bool, 1),
	}
}

func (s *worker) Start() {
	go func() {
		for {
			item, err := s.queue.ReadNext()
			if err != nil {
				if errors.Is(err, queue.ErrorQueueClosing) {
					return
				}
				log.Println(err)
				continue
			}

			cmd := s.parser.ParseCommand(item)
			s.executeCommand(cmd)

			select {
			case <-s.stop:
				return
			default:
			}
		}
	}()
}

func (w *worker) Stop() {
	w.stop <- true
}

func (w *worker) executeCommand(cmd *commandparser.Command) {
	switch cmd.Type {
	case commandparser.AddItemCommand:
		w.data.Set(cmd.Key, cmd.Value)
	case commandparser.DeleteItemCommand:
		w.data.Delete(cmd.Key)
	case commandparser.GetItemCommand:
		w.writeSingleValue(cmd.Key)
	case commandparser.GetAllItemsCommand:
		w.writeAllValues()
	}
}

func (w *worker) writeSingleValue(key string) (string, error) {
	value := w.data.Get(key)
	filename := uuid.NewString()
	path := w.outputFolder + filename
	f, err := os.Create(path)
	if err != nil {
		log.Panic(err)
		return "", err
	}
	f.WriteString(fmt.Sprintf("%s:%s", key, value))
	f.Close()
	return filename, nil
}

func (w *worker) writeAllValues() (string, error) {
	filename := uuid.NewString()
	path := w.outputFolder + filename
	f, err := os.Create(path)
	if err != nil {
		log.Panic(err)
		return "", err
	}
	w.data.Iterate(func(key, value string) {
		f.WriteString(fmt.Sprintf("%s:%s\n", key, value))
	})
	f.Close()
	return filename, nil
}
