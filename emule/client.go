package emule

import (
	"fmt"
	// "io"
	"net"
	util "github.com/AltTechTools/gomule-tst/emule"
	"time"
	//"github.com/test3-damianfurrer/gomule/tree/sharedtest/emule"
	libdeflate "github.com/4kills/go-libdeflate/v2"
)

type Client struct {
	Server     string
	Port       int
	Username   string
	Uuid	   []byte
	Debug      bool
	Ctcpport   int
	ClientConn net.Conn
	Comp	   libdeflate.Compressor
	DeComp	   libdeflate.Decompressor
	SrvTCPCompression         bool
	SrvTCPNewTags             bool
	SrvTCPUnicode             bool
	SrvTCPRelatedSearch       bool
	SrvTCPTypeTagInterger     bool
	SrvTCPLargeFiles          bool
	SrvTCPObfuscation         bool
}

func NewClientConn(server string, port int, debug bool) *Client {
	return &Client{
		Server:   server,
		Port:     port,
		Username: "gomuleclientuser",
		Ctcpport: 4662,
		Debug:   debug}
}
func (this *Client) SetTCPFlags(tcpmap uint32){
	this.SrvTCPCompression		= false
	this.SrvTCPNewTags		= false
	this.SrvTCPUnicode		= false
	this.SrvTCPRelatedSearch 	= false
	this.SrvTCPTypeTagInterger 	= false
	this.SrvTCPLargeFiles 		= false
	this.SrvTCPObfuscation 		= false
	
	if tcpmap & uint32(0x00000001) != 0 {
		this.SrvTCPCompression = true
		fmt.Println("this.SrvTCPCompression")
	}
	if tcpmap & uint32(0x00000008) != 0 {
		this.SrvTCPNewTags = true
		fmt.Println("this.SrvTCPNewTags")
	}
	if tcpmap & uint32(0x00000010) != 0 {
		this.SrvTCPUnicode = true
		fmt.Println("this.SrvTCPUnicode")
	}
	if tcpmap & uint32(0x00000040) != 0{
		this.SrvTCPRelatedSearch = true
		fmt.Println("this.SrvTCPRelatedSearch")
	}
	if tcpmap & uint32(0x00000080) != 0 {
		this.SrvTCPTypeTagInterger = true
		fmt.Println("this.SrvTCPTypeTagInterger")
	}
	if tcpmap & uint32(0x00000100) != 0 {
		this.SrvTCPLargeFiles = true
		fmt.Println("this.SrvTCPLargeFiles")
	}
	if tcpmap & uint32(0x00000400) != 0 {
		this.SrvTCPObfuscation = true
		fmt.Println("this.SrvTCPObfuscation")
	}
/*
		// Server TCP flags
#define SRV_TCPFLG_COMPRESSION          0x00000001
#define SRV_TCPFLG_NEWTAGS                      0x00000008
#define SRV_TCPFLG_UNICODE                      0x00000010
#define SRV_TCPFLG_RELATEDSEARCH        0x00000040
#define SRV_TCPFLG_TYPETAGINTEGER       0x00000080
#define SRV_TCPFLG_LARGEFILES           0x00000100
#define SRV_TCPFLG_TCPOBFUSCATION	0x00000400
		*/

}

func (this *Client) AskServerList(){
	//test ask for serverlist
	//client.Conn
	size_b:=util.UInt32ToByte(uint32(1))
	data := []byte{0xe3,size_b[0],size_b[1],size_b[2],size_b[3],0x14}
	this.ClientConn.Write(data)
}

func (this *Client) read(conn net.Conn) (buf []byte, protocol byte, err error) {
	buf = make([]byte, 5)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("ERROR: reading response:", err.Error())
		return
	}
	protocol = buf[0]
	size := util.ByteToUint32(buf[1:n])
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
		for i := 1; i < 5; i++ {
		n, err = conn.Read(tmpbuf)
		if err != nil {
			fmt.Println("ERROR: on read to buf", err.Error())
			fmt.Println("ERROR: on read to buf full", err)
			//return
		}
			break
		}
		if err != nil {
			return
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
		buf, protocol, err = this.read(this.ClientConn)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("ERROR: END Connection", err.Error())
			} else {
				fmt.Println("ERROR: error in response reading", err.Error())
				fmt.Println("ERROR: error in response readingall", err)
			}
			return
		}
		fmt.Printf("Protocol 0x%x ",protocol)
		handleServerMsg(protocol,buf,this)
	}
	return
}

func (this *Client) Connect() {
	var err error
	//libdeflate.Compressor
	this.Comp, err = libdeflate.NewCompressor()
	if err != nil {
		fmt.Println("ERROR: creating new Compressor: ", err.Error())
		return
	}
	this.DeComp, err = libdeflate.NewDecompressor()
	if err != nil {
		fmt.Println("ERROR: creating new Decompressor: ", err.Error())
		return
	}
	
	this.ClientConn, err = net.Dial("tcp",fmt.Sprintf("%s:%d",this.Server,this.Port))
	if err != nil {
		fmt.Println("ERROR: connecting: ", err.Error())
		return
	}
	body := make([]byte,0)
	//body = append(body,0x6a,0xff,0x9d,0x13,0xba,0x4f,0x4b,0x67,0xaf,0x0c,0xf6,0xa5,0x14,0xc4,0xd4,0x99) //client uuid this.Uuid
	body = append(body,this.Uuid...) //client uuid
	abuf := util.UInt32ToByte(uint32(0))
	body = append(body,abuf...) //client id 0 default
	body = append(body,util.UInt16ToByte(uint16(this.Ctcpport))...) //tcp port default
	body = append(body,util.UInt32ToByte(uint32(4))...) //tag count
	body = append(body,util.EncodeByteTagString(util.EncodeByteTagNameInt(0x1),this.Username)...)
	body = append(body,util.EncodeByteTagInt(util.EncodeByteTagNameInt(0x11),uint32(0x3C))...)
	body = append(body,util.EncodeByteTagInt(util.EncodeByteTagNameInt(0x20),uint32(0b1100011101))...)
	body = append(body,util.EncodeByteTagInt(util.EncodeByteTagNameInt(0xfb),util.ByteToUint32([]byte{128, 13, 4, 3}))...)
	//body = append(body,util.EncodeByteTagInt(util.EncodeByteTagNameInt(0x20),)...)
	
	fmt.Println("Size body", len(body))
	data := util.EncodeByteMsg(0xE3,0x01,body)
	this.ClientConn.Write(data)
	time.Sleep(10*time.Second)
	this.ConnReader() //reads all incoming data
	return
}
	
func (this *Client) Disconnect() {
	//defer this.listener.Close()
	defer this.ClientConn.Close()
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

/*

[56 56 0 69 82 82 79 82 32 58 32 89 111 117 114 32 101 100 111 110 107 101 121 32 99 108 105 101 110 116 32 105 115 32 116 111 111 32 111 108 100 44 32 112 108 101 97 115 101 32 117 112 100 97 116 101 32 105 116]
Protocol e3Received buf:  [56 83 0 87 65 82 78 73 78 71 32 58 32 89 111 117 32 104 97 118 101 32 97 32 108 111 119 105 100 46 32 80 108 101 97 115 101 32 114 101 118 105 101 119 32 121 111 117 114 32 110 101 116 119 111 114 107 32 99 111 110 102 105 103 32 97 110 100 47 111 114 32
*/
