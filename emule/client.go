package emule

import (
	"fmt"
	// "io"
	"net"
	_ "github.com/test3-damianfurrer/gomule@sharedtest/emule"
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

func (this *Client) read(conn net.Conn) (buf []byte, protocol byte, err error) {
	buf = make([]byte, 5)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("ERROR: reading response:", err.Error())
		return
	}
	protocol = buf[0]
	size := ByteToUint32(buf[1:n])
	//if this.Debug {
	//	fmt.Printf("DEBUG: size %v -> %d\n", buf[1:n], size)
	//}
	buf = make([]byte, 0)
	toread := size
	var tmpbuf []byte
	for{
		if toread > 1024  {
			tmpbuf = make([]byte, 1024)
		} else {
			tmpbuf = make([]byte, toread)
		}
		n, err = conn.Read(tmpbuf)
		if err != nil {
			fmt.Println("ERROR: on read to buf", err.Error())
			//return
		}
		buf = append(buf, tmpbuf[0:n]...)
		if n < 0 {
			fmt.Println("WARNING: n (conn.Read) < 0, some problem")
			n = 0
		}
		toread -= uint32(n)
		if toread <= 0 {
			break;
		}
	}
	return
	
}
func (this *Client) ConnReader() {
	var buf []byte
	var protocol byte
	var err error
	for {
		buf, protocol, err = this.read(this.ClientConn,)
		if err != nil {
			fmt.Println("ERROR: error in response reading", err.Error())
		}
		fmt.Printf("Protocol %x",protocol)
		fmt.Println("Received buf: ", buf)
	}
	return
}

func (this *Client) Connect() {
	var err error
	this.ClientConn, err = net.Dial("tcp",fmt.Sprintf("%s:%d",this.Server,this.Port))
	if err != nil {
		fmt.Println("ERROR: connecting: ", err.Error())
	}
	var body []byte
	body = append(body,0x6a,0xff,0x9d,0x13,0xba,0x4f,0x4b,0x67,0xaf,0x0c,0xf6,0xa5,0x14,0xc4,0xd4,0x99) //client uuid
	body = append(body,Uint32ToByte(uint32(0))...) //client id 0 default
	body = append(body,Uint16ToByte(uint16(4662))...) //tcp port default
	body = append(body,Uint32ToByte(uint32(3))...) //tag count
	body = append(body,EncodeByteTagString(encodeByteTagNameInt(0x1),"gomuleclientuser"))
	body = append(body,EncodeByteTagInt(encodeByteTagNameInt(0x11),uint32(0x3C)))
	body = append(body,EncodeByteTagInt(encodeByteTagNameInt(0x20),uint32(0b1100011101)))
	
	data := EncodeByteMsg(0xE3,0x01,body)
	this.ClientConn.Write(data)
	this.ConnReader()
	return
}
	
func (this *Client) Disconnect() {
	//defer this.listener.Close()
	return
}

/*
name tag len 1
DEBUG: len(tagarr) 4
Debug Name Tag: http://www.aMule.org1
Debug Version Tag: 60
Debug Flags Tag: 1100011101
Warning: unknown tag 0xfb
*/
