package server

import (
	"coding-exercise/internal/commandparser"
	"coding-exercise/internal/orderedmap"
	"coding-exercise/internal/queue"
)

type Server struct {
	data           orderedmap.OrderedMap
	queue          queue.Queue
	parser         commandparser.CommandParser
	outputFolder   string
	maxParallelism int
	workers        []Worker
}

func NewServer(maxParallelism int, queue queue.Queue, outputFolder string) *Server {
	return &Server{
		data:           orderedmap.NewOrdererMap(),
		queue:          queue,
		parser:         commandparser.NewCommandParser(),
		outputFolder:   outputFolder,
		maxParallelism: maxParallelism,
	}
}

func (s *Server) Run() {
	for i := 0; i < s.maxParallelism; i++ {
		worker := newWorker(s.data, s.queue, s.parser, s.outputFolder)
		s.workers = append(s.workers, worker)
		worker.Start()
	}
}

func (s *Server) Stop() {
	for _, worker := range s.workers {
		worker.Stop()
	}
}

// scale function provided as example,
// could be used from a queue monitor and scale out/in depending on the nr of messages in the queue
func (s *Server) Scale(size int) {
	if size > 0 {
		for i := 0; i < size; i++ {
			worker := newWorker(s.data, s.queue, s.parser, s.outputFolder)
			s.workers = append(s.workers, worker)
			worker.Start()
		}
	} else {
		numWorkers := len(s.workers)
		for i := numWorkers - 1; i > numWorkers+size; i-- {
			s.workers[i].Stop()
		}

		s.workers = s.workers[:numWorkers+size]
	}
}
