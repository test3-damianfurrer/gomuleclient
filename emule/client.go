package emule

import (
	"fmt"
	"io"
	"net"
)

type Client struct {
	Server   string
	Port     int
	Debug    bool
}

func NewClientConn(server string, port int, debug bool) *Client {
	return &Client{
		Server:  server,
		Port:    port,
		Debug:   debug}
}
