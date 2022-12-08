package emule

import (
	"fmt"
	"io"
	"net"
)

type Client struct {
	Server     string
	Port       int
	Debug      bool
	ClientConn net.Conn
}

func NewClientConn(server string, port int, debug bool) *Client {
	return &Client{
		Server:  server,
		Port:    port,
		Debug:   debug}
}


func (this *Client) Connect() {
	var err error
	this.ClientConn, err = net.Dial("tcp",fmt.Sprintf("%s:%d",this.Server,this.Port))
	
}
	
func (this *Client) Disconnect() {
	//defer this.listener.Close()
	return
}
