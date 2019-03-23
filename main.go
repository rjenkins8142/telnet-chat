package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/rjenkins8142/tchat/chatroom"
	"github.com/rjenkins8142/tchat/config"
	"github.com/rjenkins8142/tchat/version"
	"github.com/spf13/viper"
)

func main() {
	var configFile string
	const (
		defaultConfigFile = "config.toml"
		configUsage       = "Full path/filename of the config file"
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (%s):\n", filepath.Base(os.Args[0]), version.Info())
		flag.PrintDefaults()
	}

	flag.StringVar(&configFile, "config", defaultConfigFile, configUsage)
	flag.StringVar(&configFile, "c", defaultConfigFile, configUsage)

	flag.Parse()

	log.Printf("Reading config from %s...\n", configFile)

	config.ParseConfig(configFile)

	tcpAddr := viper.GetString("tcp.addr")
	tcpPort := viper.GetString("tcp.port")

	listener, err := net.Listen("tcp", tcpAddr+":"+tcpPort)
	if err != nil {
		log.Fatalf("Error listening on %s:%s: %s\n", tcpAddr, tcpPort, err)
	}
	defer listener.Close()

	// Create default lobby room.
	lobby, err := chatroom.CreateRoom("lobby")
	if err != nil {
		log.Fatalf("Error creating lobby: %s\n", err)
	}

	log.Printf("Listening on %s:%s\n", tcpAddr, tcpPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s\n", err)
		}
		go chatroom.CreateUser(conn, lobby)
	}
}
