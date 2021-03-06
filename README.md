# Simple telnet chat server

## Prerequisites:

Go 1.9+ (written with Go 1.12, but should work with most versions of Go that support dep.)

## To clone:

```
$ go get github.com/rjenkins8142/telnet-chat
$ cd $GOPATH/src/github.com/rjenkins8142/telnet-chat
```

## To build:

If GNU Make is installed:

`$ make`

otherwise:

`$ go build`

#### Dependency management is done with the default dep tool.

## Usage:

### Start server

`$ ./telnet-chat`

Edit the config.toml file to change any default configuration settings.

You can also use environment variables to override any config settings.

For instance:
```
TCHAT_TCP_ADDR
TCHAT_TCP_PORT
TCHAT_LOG_FILEPATH
```

### Connect via telnet

`$ telnet 127.0.0.1 8080`

## TODO

* Unit tests
* HTTP/RESTful API
* Additional Chat Commands
* Better terminal handling?

