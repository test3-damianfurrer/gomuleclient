package emule

import (
	"fmt"
)

func handleServerMsg(protocol byte,buf []byte){
    bufsize := len(buf)
    if protocol = 0xe3 {
        switch buf[0] {
            case 0x38:
            prcServerTextMsg(buf[1:bufsize])
            default:
            fmt.Printf("ERROR: Msg type %x not supported\n",buf[0])
        }
    } else {
        //decode
        fmt.Println("ERROR: only std 0xE3 protocol supported")
    }
}

func prcServerTextMsg(buf []byte){
    
}
