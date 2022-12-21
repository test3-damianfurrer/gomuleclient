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
	fmt.Println("DEBUG: decompressed",decompressed[0:40])
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
			case 0x33:
				prcSearchResults(buf)
			case 0x32:
				prcServerList(buf)
			case 0x41:
				prcServerIdentification(buf)
            default:
            	fmt.Printf("ERROR: Msg type 0x%x not supported\n",btype)
        }
}

func prcOneSearchResult(pos int, buf []byte) (readb int, fname_b []byte, hash_b []byte){
	readb=pos
	hash_b=buf[readb:readb+16]
	fmt.Printf("Debug: hash: 0x%x \n",hash_b)
	readb+=16
	fmt.Println("Debug: peer ip: ",buf[readb:readb+4])
	readb+=4
	fmt.Println("Debug: peer port: ",buf[readb:readb+2])
	readb+=2
	tagcount:=util.ByteToUint32(buf[readb:readb+4])
	fmt.Println("Debug: tag count: ",buf[readb:readb+4],tagcount)
	readb+=4
	//tag count is wrong man break.
	forbreak:=false
	for {
		//[157 1 85 45 66 111 111 116 32 57 54 46 97 118 105 131 2 0 224 21 126 137 21 11 137 48 11 136 212 124 5 148 213 100 120 53 48 136 211 204 46 26 255 154 116 172 162 235 163 143 237 166 47 133 235 224 59 228 179 14 0 82 80 4 0 0 0 130 1 45 0 68 97 115 32 66 111 111 116 32 49 120 48 53 32 45 32 69 112 105 115 111 100 105 111 32 48 53 32 91]
		//85 45 66 111 111 116 32 57 54 46 97 118 105
		//U-Boot 96.avi
		//chars: 13
		//no len encoded
		fmt.Println("Debug: tag indicator", buf[readb])
		fmt.Println("Debug: tag indicator++", buf[readb:readb+5])
		switch buf[readb] {
			case 100:
				switch buf[readb+1] {
					case 105:
						fmt.Println("Debug: some tagging/value: ",buf[readb:readb+4])
						readb+=4 //idk, what this should be
					case 120:
						if buf[readb+2] == 53 && buf[readb+3] == 48 { //dx50
							fmt.Println("Debug: dx50 tagging: ",buf[readb:readb+4])
							readb+=4
							fmt.Println("Debug: dx50 value: ",buf[readb:readb+4])
							readb+=4
						} else {
							forbreak=true
						}
					default:
						forbreak=true
				}
				
			case 120:
				if buf[readb+1] == 118 && buf[readb+2] == 105 && buf[readb+3] == 100 { //118 105 100 == vid
					fmt.Println("Debug: vid tagging: ",buf[readb:readb+4])
					readb+=4
					fmt.Println("Debug: unknown value: ",buf[readb:readb+4],util.ByteToUint32(buf[readb:readb+4]))
					readb+=4
				} else {
					//break
					forbreak=true
				}
			case 130:
				if buf[readb+1] == 1 { //filename
					fmt.Println("Debug: tagging: ",buf[readb:readb+2])
					readb+=2
					strlen := util.ByteToUint16(buf[readb:readb+2])
					readb+=2
					fmt.Println("Debug: strlen",strlen)
					fmt.Println("Debug: str",buf[readb:readb+int(strlen)])
					fname_b=buf[readb:readb+int(strlen)]
					fmt.Printf("Debug: str: %s\n",fname_b)
					readb+=int(strlen)
				} else {
					//break
					forbreak=true
				}
			case 131:
				if buf[readb+1] == 2 {
					fmt.Println("Debug: unknown tagging: ",buf[readb:readb+2])
					readb+=2
					fmt.Println("Debug: unknown value: ",buf[readb:readb+4],util.ByteToUint32(buf[readb:readb+4]))
					readb+=4
				} else {
					//break
					forbreak=true
				}
			case 136:
				switch buf[readb+1] {
					case 2:
						//136 2 13 1 //bs?
						fmt.Println("Debug: unknown tag/value: ",buf[readb:readb+4])
						readb+=4
					case 212:
						fmt.Println("Debug: unknown tagging: ",buf[readb:readb+2])
						readb+=2
						fmt.Println("Debug: unknown value: ",buf[readb:readb+4],util.ByteToUint32(buf[readb:readb+4]))
						readb+=4
					case 247: //bs
						fmt.Println("Debug: unknown tag/value: ",buf[readb:readb+4])
						readb+=4
					case 211: //bs
						fmt.Println("Debug: unknown tag/value: ",buf[readb:readb+4])
						readb+=4
					default:
						forbreak=true
				}
			case 137:
				switch buf[readb+1] {
					case 21:
						fmt.Println("Debug: unknown tagging: ",buf[readb:readb+2])
						readb+=2
						fmt.Println("Debug: unknown value: ",buf[readb:readb+4],util.ByteToUint32(buf[readb:readb+4]))
						readb+=4
					case 212: //bs
						fmt.Println("Debug: unknown tag/value: ",buf[readb:readb+3]) //only 1?
						readb+=3
					default:
						forbreak=true
				}
			case 139:
				if buf[readb+1] == 1 {
					readb+=2 //idk, what this should be
				} else {
					forbreak=true
				}
			case 156,157:
				if buf[readb+1] == 1 {
					readb+=1
					bufstr:=make([]byte,0)
					for  {
						readb+=1
						breakbread := false
						switch buf[readb] {
							case 131,136:
								breakbread=true
						}
						if breakbread {
							break
						}
						bufstr=append(bufstr,buf[readb])
					}
					fmt.Println("(obfuscated?)name buf:",bufstr)
					fmt.Printf("(obfuscated?)name buf:%s\n",bufstr)
					fname_b=bufstr
				} else {
					forbreak=true
				}
			case 213:
				readb+=1 //idk, what this should be
			default:
				forbreak=true //break
		}
		if forbreak {
			break
		}
	}
	//fmt.Println("Debug: skipped: ",buf[readb:readb+2])
	//readb+=2
	//strlen := util.ByteToUint16(buf[readb:readb+2])
	//readb+=2
	//fmt.Println("Debug: strlen",strlen)
	//fmt.Println("Debug: str",buf[readb:readb+int(strlen)])
	//readb+=int(strlen)
	//fmt.Println("Debug: tag indicator", buf[readb])
	//fmt.Println("Debug: tag indicator++", buf[readb:readb+5])
	/*
	177 214 124 66 210 36 9 186 229 94 162 137 249 109 238 248 
	81 184 64 167 
	54 18 
	7 0 0 0
	
	100 105 118 51 
	136 211 38 21 
	149 15 208 203 45 213 224 178 126 67 15 168 114 8 151 64 
	95 235 173 246 
	54 18 
	7 0 0 0
	*/
	
	readb-=pos
	return
}

func prcSearchResults(buf []byte){
	rescount := util.ByteToUint32(buf[0:4])
	fmt.Println("Debug: search rescount: ",rescount)
	
	prcread := 0
	i := 0
	bread2 := 0
	for i = 0; i<20; i++ {
		//prcread += prcOneSearchResult(4+prcread,buf)
		bread, fname_b, hash_b := prcOneSearchResult(4+prcread,buf)
		if (len(fname_b)==0){
			fmt.Println("Error: couldn't parse result file name: ",buf[4+prcread-bread2:4+prcread+bread+100])
			fmt.Println("     prev:",buf[4+prcread-bread2:4+prcread])
			fmt.Println("     curr:",buf[4+prcread:4+prcread+bread])
			fmt.Println("     +100:",buf[4+prcread+bread:4+prcread+bread+100])
			break
		}
		bread2=bread
		prcread += bread
		fmt.Println("Debug: prcread",prcread)
		fmt.Printf("\n0x%x|%s\n\n",hash_b,fname_b)
	}
	fmt.Println("Debug: after:",i,buf[4+prcread:4+prcread+100])
	//firstHash := util.ByteToUint32(buf[4:20])
	/*
	fmt.Printf("Debug: first hash: 0x%x \n",buf[4:20])
	fmt.Println("Debug: peer ip: ",buf[20:24])
	fmt.Println("Debug: peer port: ",buf[24:26])
	fmt.Println("Debug: tag count: ",buf[26:30],util.ByteToUint32(buf[26:30]))
	//fmt.Println("Debug: after: ",buf[30:60])
	fmt.Println("Debug: skipped: ",buf[30:32])
	strlen := util.ByteToUint16(buf[32:34])
	fmt.Println("Debug: strlen",strlen)
	fmt.Println("Debug: str",buf[34:34+strlen])
	///fmt.Println("Debug: after str",buf[34+strlen:34+strlen+120])
	fmt.Println("Debug: skipped: ",buf[34+strlen:34+strlen+2])
	fmt.Println("Debug: skipped val: ",buf[34+strlen+2:34+strlen+6])
	fmt.Println("Debug: skipped: ",buf[34+strlen+6:34+strlen+8])
	iend:=34+strlen+12
	fmt.Println("Debug: skipped val: ",buf[34+strlen+8:iend])
	
	prcread := 0
	prcread += prcOneSearchResult(4+prcread,buf)
	fmt.Println("Debug: prcread",prcread)
	//fmt.Println("Debug: (maybe + 12)should be",34+strlen-4)
	//iend:=4+prcread
	prcread += prcOneSearchResult(4+prcread,buf)
	fmt.Println("Debug: prcread",prcread)
	prcread += prcOneSearchResult(4+prcread,buf)
	fmt.Println("Debug: prcread",prcread)
	fmt.Println("Debug: after 3:",buf[4+prcread:4+prcread+100])
	
	fmt.Printf("Debug: second hash: 0x%x \n",buf[iend:iend+16])
	fmt.Println("Debug: second ip: ",buf[iend+16:iend+20])
	fmt.Println("Debug: second port: ",buf[iend+20:iend+22])
	fmt.Println("Debug: second tag count: ",buf[iend+22:iend+26])
	
	fmt.Println("Debug: skipped: ",buf[iend+26:iend+28])
	strlen2 := int(util.ByteToUint16(buf[iend+28:iend+30]))
	fmt.Println("Debug: strlen2",strlen2)
	fmt.Println("Debug: str2",buf[iend+30:iend+30+strlen2])
	fmt.Println("Debug:second after: ",buf[iend+30+strlen2:iend+30+strlen2+100])
	*/
	//248 1 0 0 -> 504 results
	
	// [201 184 52 216 95 73 187 94 17 15 11 174 35 74 120 95] 
	// [50 53 142 41]
	// [70 57]
	//4 0 0 0  - tag count 
	//130 1 - fname tag indicator
	// --- Linux.Magazine.164.Jul.2014--Cygwin, UEFI Secure Boot, HDD to SSD, Bitwig Studio, Bash History, Oculus Rift VR.pdf
	//131 2 28 136 91 1 137 21 7 137 48 7 
	//94 20 60 87 136 226 150 127 40 123 187 107 79 243 48 29
	//151 83 189 58
	//229 6 
	//7 0 0 0 130 1 
	// --- Tae Bo - Billys' Boot Camp - Basic Training - cd 1 - (Billy Blanks) - [ENG].avi
	//131 2 0 66 227 43
	//137 21 1 137 48 1 
	//136 212 1 14 148 213 120 118 105 100 136 211 102 6 
	//40 240 189 252 88 72 242 10 172 142 231 225 73 69 107 219 
	//93 186 255 106 
	//54 18 
	//4 0 0 0 130 1
	//Flavors of Puglia. Traditional recipes from.... 
	
	//hashs:
	//201 184 52 216 95 73 187 94 17 15 11 174 35 74 120 95
	//94 20 60 87 136 226 150 127 40 123 187 107 79 243 48 29
	//40 240 189 252 88 72 242 10 172 142 231 225 73 69 107 219
	
	//ips:
	//50 53 142 41
 	//151 83 189 58
	//93 186 255 106 
	
	//ports:
	//70 57
	//229 6
	//54 18
	
	//tag counts: (might be intentionally wrong?)
	//4 0 0 0
	//7 0 0 0
	//4 0 0 0
	
	//tags
	//131 2 + 4bytes
	//137 21 + 4bytes
	
	//136 212 + 4bytes
	//120 118 105 100 + 4bytes (string tag name "vid")
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
	client.SearchServer()
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
