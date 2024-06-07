package main

import (
	"bufio"
	"coding-exercise/internal/queue"
	"coding-exercise/internal/server"
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "rabbitmq_endpoint",
				Value: "amqp://guest:guest@localhost:5672/",
				Usage: "the url of the rabbitmq server to be used",
			},
			&cli.StringFlag{
				Name:  "rabbitmq_queue_name",
				Value: "testq",
				Usage: "the queuename of the queue in rabbitmq",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "consumer",
				Usage: "starts a consumer",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "outputfolder",
						Value: "./out/",
						Usage: "the folder where to write the output of the read commands",
					},
					&cli.IntFlag{
						Name:  "concurrency",
						Value: 10,
						Usage: "the maximum concurrency of the consumer",
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					q := createQueue(cmd)
					defer q.Close()
					concurrency := cmd.Int("concurrency")
					outputFolder := cmd.String("outputfolder")
					startConsumer(q, outputFolder, int(concurrency))
					return nil
				},
			},
			{
				Name:  "producer",
				Usage: "starts a producer",
				Commands: []*cli.Command{
					{
						Name:            "command",
						Usage:           "produce a single command from command line",
						HideHelpCommand: true,
						ArgsUsage:       "<command>",
						Action: func(_ context.Context, cmd *cli.Command) error {
							q := createQueue(cmd)
							defer q.Close()
							command := cmd.Args().First()
							if command == "" {
								return errors.New("specify a command")
							}
							return q.Publish(command, time.Second*5)
						},
					},
					{
						Name:            "file",
						Usage:           "starts a producer",
						HideHelpCommand: true,
						ArgsUsage:       "<filename>",
						Action: func(_ context.Context, cmd *cli.Command) error {
							q := createQueue(cmd)
							defer q.Close()
							commandFile := cmd.Args().Get(0)
							if commandFile == "" {
								return errors.New("specify a filename")
							}
							runFileProducer(q, commandFile)
							return nil
						},
					},
					{
						Name:  "random",
						Usage: "starts a random producer",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "read",
								Value: true,
								Usage: "produce read commands",
							},
							&cli.BoolFlag{
								Name:  "write",
								Value: true,
								Usage: "produce write commands",
							},
						},
						Action: func(_ context.Context, cmd *cli.Command) error {
							q := createQueue(cmd)
							defer q.Close()
							read := cmd.Bool("read")
							write := cmd.Bool("write")
							runRandomProducer(q, read, write)
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func createQueue(cmd *cli.Command) queue.Queue {
	url := cmd.String("rabbitmq_endpoint")
	name := cmd.String("rabbitmq_queue_name")
	q := queue.NewRabbitMQ(url, name)
	err := q.Open()
	if err != nil {
		panic(err)
	}
	return q
}

func waitForSigTerm() {
	log.Println("Press CTRL+C to exit")
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	<-cancelChan
}

func startConsumer(q queue.Queue, outputFolder string, concurrency int) {
	log.Println("Starting consumer")
	q.OpenConsumer()
	if err := os.MkdirAll(outputFolder, 0777); err != nil {
		log.Panic(err)
	}

	srv := server.NewServer(concurrency, q, outputFolder)
	srv.Run()
	waitForSigTerm()
	srv.Stop()
}

func runFileProducer(q queue.Queue, filename string) {
	log.Println("Publishing commands from file")
	file, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}

	rdr := bufio.NewReader(file)
	for line, isPrefix, err := rdr.ReadLine(); err == nil; line, isPrefix, err = rdr.ReadLine() {
		if !isPrefix && len(line) > 0 {
			log.Println(string(line))
			q.Publish(string(line), time.Second*5)
		}
	}
}

func runRandomProducer(q queue.Queue, read, write bool) {
	log.Println("Starting producer")
	log.Println("Press CTRL+C to exit")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	imq := queue.NewInMemeoryQueue(read, write)
	for {
		msg, _ := imq.ReadNext()
		q.Publish(msg, time.Second*5)

		select {
		case <-cancelChan:
			return
		default:
		}
	}
}
