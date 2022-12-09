package emule

import (
	"fmt"
	util "github.com/AltTechTools/gomule-tst/emule"
)

func handleServerMsg(protocol byte,buf []byte){
    bufsize := len(buf)
	if protocol == 0xe3 {
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
    } else {
        //decode
        fmt.Println("ERROR: only std 0xE3 protocol supported")
    }
}

func prcServerIdentification(buf []byte){
	serveruuid:=buf[0:16]
	serverip:=util.ByteToUint32(buf[16:20])
	serverport:=util.ByteToUint32(buf[20:24])
	tags:=util.ByteToUint32(buf[24:28])
	
	fmt.Printf("server uuid %x-%x-%x-%x\n",serveruuid[0:4],serveruuid[4:8],serveruuid[8:12],serveruuid[12:16])
	fmt.Println("server ip",serverip)
	fmt.Println("server port",serverport)
	fmt.Println("msg tags",tags)
	
	fmt.Println("bytes tags",buf[28:len(buf)])
	
	
}

func prcServerStatus(buf []byte){
	usercount:=util.ByteToUint16(buf[0:4])
	filecount:=util.ByteToUint16(buf[4:8])
	fmt.Println("Server Users",usercount)
	fmt.Println("Server Files",filecount)
}

func prcIdChange(buf []byte){
	clientid:=util.ByteToUint16(buf[0:4])
	fmt.Println("Client id",clientid)
	if len(buf) == 8 {
		tcpmap:=util.ByteToUint16(buf[4:8])
		fmt.Printf("tcp map %b\n",tcpmap)
	}
}

func prcServerTextMsg(buf []byte){
	msglen := util.ByteToUint16(buf[0:2])
	fmt.Printf("String: \n%s\n",buf[2:msglen+2])
	//util.readString(0,buf)
}
