package emule

import (
	"fmt"
	"io"
	"net"

	sam "github.com/eyedeekay/sam3/helper"
)

type Client struct {
	Server   string
	Port     int
	Debug    bool
}

func NewClient(server string, port int, debug bool) *Client {
	return &SockSrv{
		Server:  server,
		Port:    port,
		Debug:   debug}
}
