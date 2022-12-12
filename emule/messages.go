package emule

import (
	"fmt"
	util "github.com/AltTechTools/gomule-tst/emule"
	libdeflate "github.com/4kills/go-libdeflate/v2" //libdeflate.Compressor
)

func handleServerMsg(protocol byte,buf []byte, client *Client){
    	//0xd4
	switch protocol {
		case 0xe3:
			decodeE3(buf[0],buf[1:len(buf)],client)
		case 0xd4:
			decodeD4(buf[0],buf[1:len(buf)],client.DeComp,client)
		default:
			fmt.Println("ERROR: only std 0xE3 protocol supported")
	}
}

func decodeD4(btype byte,buf []byte,dc libdeflate.Decompressor, client *Client){
	fmt.Printf("DEBUG: 0xd4 type 0x%x\n",btype)
	blen, decompressed, err := dc.Decompress(buf, nil, 1)
	if err != nil {
		fmt.Println("ERROR: failed to decompress buffer",err)
		return
	}
	fmt.Println("DEBUG: decompressed length:",blen)
	fmt.Println("DEBUG: decompressed",decompressed[0:30])
	decodeE3(btype,decompressed,client)
}

func decodeE3(btype byte,buf []byte, client *Client){
	switch btype {
			case 0x38:
				prcServerTextMsg(buf)
			case 0x40:
				prcIdChange(buf,client)
			case 0x34:
				prcServerStatus(buf)
			case 0x32:
				prcServerList(buf)
			case 0x41:
				prcServerIdentification(buf)
            default:
            	fmt.Printf("ERROR: Msg type 0x%x not supported\n",btype)
        }
}

func prcServerList(buf []byte){
	offset := 1
	for i := byte(0); i < buf[0]; i++ {
		fmt.Printf("Server ip: %d.%d.%d.%d:",buf[offset],buf[offset+1],buf[offset+2],buf[offset+3])
		offset+=6
		fmt.Printf("%d\n",util.ByteToUint16(buf[offset-2:offset]))
	}
}

func prcServerIdentification(buf []byte){
	fmt.Println("Server Identification")
	serveruuid:=buf[0:16]
	serverip:=util.ByteToUint32(buf[16:20])
	serverport:=util.ByteToUint16(buf[20:24])
	tags:=util.ByteToUint32(buf[22:26])
	
	fmt.Printf("server uuid %x-%x-%x-%x\n",serveruuid[0:4],serveruuid[4:8],serveruuid[8:12],serveruuid[12:16])
	fmt.Println("server ip",serverip)
	fmt.Printf("server ip %d.%d.%d.%d\n",buf[16],buf[17],buf[18],buf[19])
	fmt.Println("server port",serverport)
	fmt.Println("msg tags",tags)
	
	fmt.Println("bytes tags",buf[26:len(buf)])
	
	nlen := int(util.ByteToUint16(buf[30:32]))
	fmt.Printf("s1 %s\n",buf[32:32+nlen])
	nlen2 := int(util.ByteToUint16(buf[32+4+nlen:32+4+nlen+2]))
	fmt.Printf("s2 %s\n",buf[32+4+nlen+2:32+4+nlen+2+nlen2])
	
	
}

func prcServerStatus(buf []byte){
	fmt.Println("Server Status")
	usercount:=util.ByteToUint32(buf[0:4])
	filecount:=util.ByteToUint32(buf[4:8])
	fmt.Println("Server Users",usercount)
	fmt.Println("Server Files",filecount)
}

func prcIdChange(buf []byte, client *Client){
	fmt.Println("ID change")
	fmt.Println("ID change tcp flags schould contain support for large files indicator")
	fmt.Printf("ID change len: %d",len(buf))
	fmt.Println(", buf: ", buf)
	clientid:=util.ByteToUint32(buf[0:4])
	fmt.Println("Client id",clientid)
	if len(buf) == 8 {
		tcpmap:=util.ByteToUint32(buf[4:8])
		fmt.Printf("tcp map %b\n",tcpmap)
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
		client.SetTCPFlags(tcpmap)
		
	}
	if len(buf) == 16 {
		tcpmap:=util.ByteToUint32(buf[4:8])
		fmt.Printf("tcp map %b\n",tcpmap)
		client.SetTCPFlags(tcpmap)
		fmt.Printf("something(server port) %d",util.ByteToUint32(buf[8:12]))
		fmt.Println(" ",buf[8:12])
		fmt.Printf("something 2(my ip) %d\n",util.ByteToUint32(buf[12:16]))
		fmt.Println(" ",buf[12:16])
	}
	client.AskServerList()
	/*
	//test ask for serverlist
	//client.Conn
	size_b:=util.UInt32ToByte(uint32(1))
	data := []byte{0xe3,size_b[0],size_b[1],size_b[2],size_b[3],0x14}
	client.ClientConn.Write(data)
	*/
}

func prcServerTextMsg(buf []byte){
	fmt.Println("Server Message")
	msglen := util.ByteToUint16(buf[0:2])
	fmt.Printf("String: \n%s\n",buf[2:msglen+2])
	//util.readString(0,buf)
}
