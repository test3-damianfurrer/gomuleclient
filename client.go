package main

import (
	"flag"
	"fmt"
	"encoding/hex"
	//"github.com/test3-damianfurrer/gomule/emule"
	"github.com/test3-damianfurrer/gomuleclient/emule"
)

var (
	debug    bool
	server   string
	username string
	uuid     string
	port     int
	lport	 int
)

func init() {
	flag.BoolVar(&debug, "d", false, "Debug")
	flag.StringVar(&server, "h", "localhost", "Server address")
	flag.StringVar(&username, "u", "gomuleuser", "Username")
	flag.StringVar(&uuid, "x", "6aff9d13ba4f4b67af0cf6a514c4d499", "User UUID hex format")
	flag.IntVar(&port, "p", 7111, "Server Port number")
	flag.IntVar(&lport,"l",5662,"Tcp Client listening port")
}

func main() {
	flag.Parse()

	fmt.Println("GoMule Client Version 1.0")
	
	client := emule.NewClientConn(server, port, debug)
	client.Username=username
	uuid_b, err := hex.DecodeString(uuid)
	if err !=  nil {
		fmt.Println("provide valid hex")
		panic(err)
		
	}
	client.Uuid=uuid_b
	client.Ctcpport=lport
	//0x6a,0xff,0x9d,0x13,0xba,0x4f,0x4b,0x67,0xaf,0x0c,0xf6,0xa5,0x14,0xc4,0xd4,0x99) //client uuid this.Uuid
	client.Connect()
	defer client.Disconnect()
}
