package main

import (
	"log"
	"bytes"
	"time"
)

func xmodemRead() (success int, failed int, data bytes.Buffer){

	oBuffer := make([]byte, 1)
	dBuffer := make([]byte, 1024)

	//log.Println("prepare")
	var protocmd []byte
	protocmd = append(protocmd, NAK)

	success = 0
	failed  = 0

	var getBytes = true
	for getBytes {

		// init tranafer
		if _, err := serialPort.Write(protocmd); err != nil {
			log.Println(err)
			break
		}

		if protocmd[0] == EOT || protocmd[0] == EOF || protocmd[0] == CAN {
			log.Printf("tranfer end.")
			break
		}

		if _, err := serialPort.Read(oBuffer); err != nil {
			log.Println(err)
			break
		}

		//start receiving blocks
		if getBytes {
			myPacket := xblock{}
			bytesReceived := 0
			blockReceived := false
			for !blockReceived {
				time.Sleep(time.Millisecond * 25)
				n, err := serialPort.Read(dBuffer)
				bytesReceived = bytesReceived + n
				if err != nil {
					log.Println("Read failed:", err)
				}

				if bytesReceived >= 131 {
					myPacket.proto = oBuffer
					myPacket.packetNum = int(dBuffer[0])
					myPacket.packetInv = int(dBuffer[1])
					myPacket.payload = dBuffer[2:130]
					myPacket.checksumm = int(dBuffer[130])

					CHK := int(checksum(myPacket.payload, 0))
					if CHK == myPacket.checksumm && myPacket.checkPaylod() {
						//packet OK
						log.Printf("Checksum OK for Packet: %d\n", myPacket.packetNum)
						protocmd[0] = ACK
						success++
						data.Write(myPacket.payload)
					} else {
						//something went wrong
						if !myPacket.checkPaylod() && failed < 10 {

							if byte(myPacket.packetNum) == EOF || byte(myPacket.packetNum) == EOT {
								//EOT & EOF are no failures
								failed--
							} else {
								//message for sender
								failed++
								blockReceived = true
								protocmd[0] = NAK

							}
						}
						//stop transfer
						//log.Printf("Failed Packet (%d)\n len: %d\nData: %X\n", myPacket.packetNum, bytesReceived, dBuffer[:bytesReceived])
						failed-- //the last packet checksum must have missmatched - no error!
						protocmd[0] = CAN
						getBytes = false
					}
					blockReceived = true
				}
			}
		}
	}
	return success, failed, data
}


func xmodemSend( p []xblock) {
	oBuffer := make([]byte, 1)
	failure := 0
	success := 0
	//log.Printf("start sending %d Packets of %d bytes payload\n", len(p), len(p[0].payload))
	for _, sp := range p {
		var reSend = true
		for reSend {
			//log.Printf("send Packet: %d\n", sp.packetNum)
			sendPacket(sp)
			if _, err := serialPort.Read(oBuffer); err != nil {
				log.Println(err)
			} else {
				switch oBuffer[0] {
				case NAK: // NAK
					//receiver ask for retransmission of this block
					log.Printf("resend Packet %d\n", sp.packetNum)
					reSend = true
					failure++
				case ACK: // ACK
					//receiver accepted this block
					reSend = false
					success++
				case CAN: // CAN
					//receiver wants to quit session
					log.Printf("receiver aborted transmission at Packet %d\n", sp.packetNum)
					reSend = false
					failure++
				default:
					//should not happen
					log.Printf("unexspected answer(0x%X) for packet %d\n", oBuffer[0], sp.packetNum)
					reSend = false
				}
			}
		}
		//when receiver sends CAN - stop transmitting
		if oBuffer[0] == CAN {
			break
		}
	}
	log.Printf("upload done - Success: %d - Failures: %d\n", success, failure)

	//send EOT byte
	var eot []byte
	var err error
	eot = append(eot, EOT)
	serialPort.Write(eot)
	n := 0
	for n == 1 {
		if n, err = serialPort.Read(oBuffer); err != nil {
			log.Println(err)
		}
		if oBuffer[0] != ACK {
			log.Printf("nexpectedanswer to EOT: 0x%X\n", oBuffer[0])
		} else {
			log.Println("end of transfer")

		}
	}
}