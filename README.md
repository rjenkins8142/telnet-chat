# Simple telnet chat server

## Prerequisites:

Go 1.9+ (written with Go 1.12, but should work with most versions of Go that support dep.)

## To build:

If GNU Make is installed:

`$ make`

otherwise:

`$ go build`

## Usage:

`$ ./telnet-chat -h`

Edit the config.toml file to change any default configuration settings.

You can also use environment variables to override any config settings.

For instance:
```
TCHAT_TCP_ADDR
TCHAT_TCP_PORT
TCHAT_LOG_FILEPATH
```

#### Dependency management is done with the default dep tool.

