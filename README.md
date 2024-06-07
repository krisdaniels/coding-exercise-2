# Prerequisites
This implementation will need a local running RabbitMQ, 
you can start a default one using the provided docker-compose.yml file in the rabbitmq directory

```bash
docker composes up
```

# Usage of the command line client

```
NAME:
   main - A new cli application

USAGE:
   main [global options] [command [command options]] [arguments...]

COMMANDS:
   consumer  starts a consumer
   producer  starts a producer
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --rabbitmq_endpoint value    the url of the rabbitmq server to be used (default: "amqp://guest:guest@localhost:5672/")
   --rabbitmq_queue_name value  the queuename of the queue in rabbitmq (default: "testq")
   --help, -h                   show help (default: false)
```


# Starting the server (consumer)

```bash
go run cmd/exercise/main.go consumer
```

# Starting the client (producer)

The client can be run in 4 different ways, 
2 random producers have been provided, one will generate read commands and the other write commands,
you can run each separately or together in the same process by specifying the command line options

```bash
go run cmd/exercise/main.go producer random
```

The other 2 modes are:

read commands from a file:
```bash
go run cmd/exercise/main.go producer file commands.txt
```

spcify a command from the command line:
```bash
go run cmd/exercise/main.go producer command "getAllItems()"
```

# Building a docker container

A build script is provided in the build folder.

```
cd build
./build-docker.sh
```

# Docker example scripts

To run the docker example scripts, a local instance of rabbit mq will need ot be running, 
to start one you can use the provided docker compose script in the rabbitmq directory.

Provided scripts and example use:

- build-docker.sh: build the docker container, needs to be run first before trying the other scripts
- start-consumer.sh: start the consumer, will create and mount an out folder in the current directory and will use that for file output from the consumer, takes no params
- start-producer-random.sh: start a test producer that will randomly generate read and write operations, takes no params
- start-producer-file.sh: start a producer that reads from a file, takes the name of a file in the current directory as parameter, ex `./start-producer-file.sh commands.txt`
- start-producer-command.sh: publishes a single command to the queue, ex. `./start-producer-command.sh "getAllItems()"`
- test.sh: run unit tests and test coverage reports
