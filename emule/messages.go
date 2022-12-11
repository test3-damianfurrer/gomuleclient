package emule

import (
	"fmt"
	util "github.com/AltTechTools/gomule-tst/emule"
	libdeflate "github.com/4kills/go-libdeflate/v2" //libdeflate.Compressor
)

func handleServerMsg(protocol byte,buf []byte,dc libdeflate.Decompressor){
    	//0xd4
	switch protocol {
		case 0xe3:
			decodeE3(buf)
		case 0xd4:
			fmt.Printf("DEBUG: maybe type 0x%x\n",buf[0])
			blen, decompressed, err := dc.Decompress(buf[1:len(buf)], nil, 1)
			if err != nil {
				fmt.Println("ERROR: failed to decompress buffer",err)
				return
			}
			fmt.Printf("DEBUG: decompressed type 0x%x\n",decompressed[0])
			fmt.Println("DEBUG: decompressed length:",blen)
			fmt.Println("DEBUG: decompressed",decompressed[0:30])
		default:
			fmt.Println("ERROR: only std 0xE3 protocol supported")
	}
}

func decodeE3(buf []byte){
	bufsize := len(buf)
	switch buf[0] {
			case 0x38:
				prcServerTextMsg(buf[1:bufsize])
			case 0x40:
				prcIdChange(buf[1:bufsize])
			case 0x34:
				prcServerStatus(buf[1:bufsize])
			case 0x41:
				prcServerIdentification(buf[1:bufsize])
            default:
            	fmt.Printf("ERROR: Msg type 0x%x not supported\n",buf[0])
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

func prcIdChange(buf []byte){
	fmt.Println("ID change")
	clientid:=util.ByteToUint32(buf[0:4])
	fmt.Println("Client id",clientid)
	if len(buf) == 8 {
		tcpmap:=util.ByteToUint32(buf[4:8])
		fmt.Printf("tcp map %b\n",tcpmap)
	}
}

func prcServerTextMsg(buf []byte){
	fmt.Println("Server Message")
	msglen := util.ByteToUint16(buf[0:2])
	fmt.Printf("String: \n%s\n",buf[2:msglen+2])
	//util.readString(0,buf)
}
