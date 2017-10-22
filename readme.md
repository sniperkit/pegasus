# Pegasus

<<<<<<< HEAD
[![Build Status](https://travis-ci.org/cpapidas/pegasus.svg?branch=master)](https://travis-ci.org/cpapidas/pegasus)
=======
[![Build Status](https://travis-ci.org/cpapidas/peg.svg?branch=master&maxAge=0)](https://travis-ci.org/cpapidas/pegasus)
>>>>>>> f85605dd5918b6a583baf12018d3f3db181fe8cb
[![Go Report Card](https://goreportcard.com/badge/github.com/cpapidas/pegasus?new=report?maxAge=0)](https://goreportcard.com/report/github.com/cpapidas/pegasus)
[![codebeat badge](https://codebeat.co/badges/d81fe30e-f110-49f1-a475-f24f1016c4c8?maxAge=0)](https://codebeat.co/projects/github-com-cpapidas-pegasus-master)
[![codecov](https://codecov.io/gh/cpapidas/pegasus/branch/master/graph/badge.svg?maxAge=0)](https://codecov.io/gh/cpapidas/pegasus)

Pegasus is an abstract layer for GRPC, AMQP (RabbitMQ) and HTTP protocols, therefore the library gives you the ability
to use all the above protocols easily. In micro-services world there are times that we need to use GRPC protocol
for internal calls, HTTP to communicate with third party services or even to expose our services and RabbitMQ
as a message broker. Pegasus helps us to do all the above.

# What Pegasus can do:

* Create HTTP server
* Create GRPC server
* Listen to RabbitMQ server
* Send requests via HTTP
* Send requests via GRPC
* Send messages to RabbitMQ server

# Get started

If you don't already set up a golang directory or GOPATH please follow the instructions bellow
[install go](https://golang.org/doc/install)
[set gopath](https://github.com/golang/go/wiki/Setting-GOPATH)

Get the project

```bash
go get github.com/cpapidas/pegasus
```

There is a `examples` folder at root directory.

In order to run a sample you have to:

```bash
$ cd examples
$ go build
```

Run the HTTP-GRPC

```bash
// run the server
$ ./examples grpchttp server
// open a new terminal and run the client
$ ./examples grpchttp client
```

*In case that you want to use RabbitMQ you have to set up and run locally a server
first.* [RabbitMQ Docker sample](https://github.com/dockerfile/rabbitmq)

Next you have to edit the *examples/sample_grpc_http_amqp/server.go* and *examples/sample_grpc_http_amqp/client.go* files
and change the rabbitMQAddress variables at top. Add your port and username/password. Usually are the same.

```bash
// run the server
$ ./examples grpchttpamqp server
// open a new terminal and run the client
$ ./examples grpchttpamqp client
```




