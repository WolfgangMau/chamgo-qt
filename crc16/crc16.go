package main

import (
	"log"
)

func main() {
	//00000000 00003230 30313138 0002031b
	//1412020f 0c0c0010 00000835 00016503
	hexstr:="10539A23"
	strhex := Split2Hex(hexstr)
	t := MakeTable(CRC16_CRC_A)
	checksum5:=Checksum(strhex, t)
	var h, l uint8 = uint8(checksum5>>8), uint8(checksum5&0xff)
	log.Printf("source: %s  XOR: %02X  crc16: %02X%02X  (%04X)\n",hexstr, GetBCC(Split2Hex(hexstr)),l,h,checksum5)
	//log.Printf("BCC: %02X : crc: %02X\n",Split2Hex("10539a25"),GetBCC(Split2Hex("0c0c00")))

}