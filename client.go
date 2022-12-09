package main

import (
	"flag"
	"fmt"
	//"github.com/test3-damianfurrer/gomule/emule"
	"github.com/test3-damianfurrer/gomuleclient/emule"
)

var (
	debug    bool
	server   string
	username string
	port     int
)

func init() {
	flag.BoolVar(&debug, "d", false, "Debug")
	flag.StringVar(&server, "h", "localhost", "Server address")
	flag.StringVar(&username, "u", "gomuleuser", "Username")
	flag.IntVar(&port, "p", 7111, "Server Port number")
}

func main() {
	flag.Parse()

	fmt.Println("GoMule Client Version 1.0")
	
	client := emule.NewClientConn(server, port, debug)
	client.Username=username
	client.Connect()
	defer client.Disconnect()
}
