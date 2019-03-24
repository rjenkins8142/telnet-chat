package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/rjenkins8142/telnet-chat/chatroom"
	"github.com/rjenkins8142/telnet-chat/config"
	"github.com/rjenkins8142/telnet-chat/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var logFile *os.File
	const (
		defaultConfigFile = "config.toml"
		configUsage       = "Full path/filename of the config file"
	)

	pflagUsage := pflag.Usage
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s (%s):\n", filepath.Base(os.Args[0]), version.Info())
		pflagUsage()
	}

	configFile := pflag.StringP("config", "c", defaultConfigFile, configUsage)

	pflag.Parse()

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	fmt.Printf("Reading config from %s...\n", *configFile)

	config.ParseConfig(*configFile)

	logPath := viper.GetString("log.filepath")

	// Handle a special case for STDOUT / STDERR for logs
	if logPath == "STDOUT" {
		logFile = os.Stdout
	} else if logPath == "STDERR" {
		logFile = os.Stderr
	} else {
		var ferr error
		// Open a new file.
		logFile, ferr = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if ferr != nil {
			log.Fatalf("Unable to open log file for writing: %s", ferr)
		}
	}

	logLevel, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		log.Fatalf("Unable to parse log Level: %s", err)
	}

	fmt.Printf("Writing logs to %s...\n", logPath)

	// Setup logging output
	log.SetOutput(logFile)
	log.SetLevel(logLevel)

	tcpAddr := viper.GetString("tcp.addr")
	tcpPort := viper.GetString("tcp.port")

	listener, err := net.Listen("tcp", tcpAddr+":"+tcpPort)
	if err != nil {
		log.Fatalf("Error listening on %s:%s: %s", tcpAddr, tcpPort, err)
	}
	defer listener.Close()

	log.Printf("---> Starting chat server <---")

	// Create default lobby room.
	lobby, err := chatroom.CreateRoom("lobby")
	if err != nil {
		log.Fatalf("Error creating lobby: %s", err)
	}
	// Initialize all the chatroom "slash" commands.
	chatroom.InitCommands()

	log.Printf("Listening on %s:%s", tcpAddr, tcpPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s", err)
		}
		go chatroom.CreateUser(conn, lobby)
	}
}
