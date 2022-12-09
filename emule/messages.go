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
            default:
            	fmt.Printf("ERROR: Msg type 0x%x not supported\n",buf[0])
        }
    } else {
        //decode
        fmt.Println("ERROR: only std 0xE3 protocol supported")
    }
}

func prcServerTextMsg(buf []byte){
	msglen := util.ByteToUint16(buf[0:2])
	fmt.Printf("String: \n%s\n",buf[2:msglen+2])
	//util.readString(0,buf)
}
