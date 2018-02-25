package xmodem

import (
	"bytes"
	"go.bug.st/serial.v1"
	"log"
	"time"
)

//xmodem
const SOH byte = 0x01
const STX byte = 0x02
const EOT byte = 0x04
const EOF byte = 0x1a
const ACK byte = 0x06
const NAK byte = 0x15
const CAN byte = 0x18

type Xblock struct {
	Proto     []byte // 1 byte protocol (SOH / STX)
	PacketNum int    // 1 byte current Packet number
	PacketInv int    // 1 byte (0xff-packetNum)
	Payload   []byte // 128 byte payload
	Checksum  int    // 1 byte complement checksum of the payload
}

func Receive(serialPort serial.Port, size int) (success int, failed int, data bytes.Buffer) {

	oBuffer := make([]byte, 1)
	dBuffer := make([]byte, size)

	//log.Println("prepare")
	var protocmd []byte
	protocmd = append(protocmd, NAK)

	success = 0
	failed = 0

	var getBytes = true
	for data.Len() < size {
		myPacket := Xblock{}

		// init tranafer
		if _, err := serialPort.Write(protocmd); err != nil {
			log.Println(err)
			break
		}

		if protocmd[0] == EOT || protocmd[0] == EOF || protocmd[0] == CAN {
			log.Printf("tranfer end.")
			getBytes=false
			if success == 0 {
				success++
				data.Write(myPacket.Payload)
			}
			break
		}

		if _, err := serialPort.Read(oBuffer); err != nil {
			log.Println(err)
			break
		}

		//start receiving blocks
		if getBytes {
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
					myPacket.Proto = oBuffer
					myPacket.PacketNum = int(dBuffer[0])
					myPacket.PacketInv = int(dBuffer[1])
					myPacket.Payload = dBuffer[2:130]
					myPacket.Checksum = int(dBuffer[130])

					CHK := int(Checksum(myPacket.Payload, 0))
					if CHK == myPacket.Checksum && myPacket.checkPaylod() {
						//packet OK
						log.Println("Checksum OK for Packet: ", myPacket.PacketNum)
						protocmd[0] = ACK
						success++
						data.Write(myPacket.Payload)
					} else {
						//something went wrong
						log.Println("something went wront with Packet: ", myPacket.PacketNum)
						if !myPacket.checkPaylod() && failed < 10 {

							if byte(myPacket.PacketNum) == EOF || byte(myPacket.PacketNum) == EOT {
								log.Println("transmission end at Block : ", myPacket.PacketNum)
								//EOT & EOF are no failures
								failed--
							} else {
								//message for sender
								log.Println("resend ... Block ", myPacket.PacketNum)
								failed++
								blockReceived = true
								protocmd[0] = NAK

							}
						}
						//stop transfer
						log.Printf("Failed Packet (%d)\n len: %d\nData: %X\n", myPacket.PacketNum, bytesReceived, dBuffer[:bytesReceived])
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

func Send(serialPort serial.Port, p []Xblock) {
	oBuffer := make([]byte, 1)
	failure := 0
	success := 0
	//log.Printf("start sending %d Packets of %d bytes payload\n", len(p), len(p[0].payload))
	for _, sp := range p {
		var resend = true //init - we need to read at least once
		for resend {
			//log.Printf("send Packet: %d\n", sp.packetNum)
			sendPacket(serialPort, sp)
			if _, err := serialPort.Read(oBuffer); err != nil {
				log.Println(err)
			} else {
				switch oBuffer[0] {
				case NAK: // NAK
					//receiver ask for retransmission of this block
					log.Printf("resend Packet %d\n", sp.PacketNum)
					resend = true
					failure++
				case ACK: // ACK
					//receiver accepted this block
					resend = false // packet was accepted, no need to resend the packet
					success++
				case CAN: // CAN
					//receiver wants to quit session
					log.Printf("receiver aborted transmission at Packet %d\n", sp.PacketNum)
					resend = false //quit session, no need to resend the packet
					failure++
				default:
					//should not happen
					log.Printf("unexspected answer(0x%X) for packet %d\n", oBuffer[0], sp.PacketNum)
					resend = true // better to read the packet again, hopefully no endless loop ;-)
					failure++
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
	eot = append(eot, CAN)
	serialPort.Write(eot)
	n := 0
	for n == 1 {
		if n, err = serialPort.Read(oBuffer); err != nil {
			log.Println(err)
		}
		if oBuffer[0] != ACK {
			log.Printf("unexpected answer to EOT: 0x%X\n", oBuffer[0])
		} else {
			log.Println("end of transfer")

		}
	}
}

func Checksum(b []byte, cs byte) byte {
	for _, d := range b {
		cs = cs + d
	}
	return cs
}

//returns false if all payload-bytes are set to 0xff
func (p Xblock) checkPaylod() bool {
	var counter = 0
	for _, b := range p.Payload {
		if b == 0xff {
			counter++
		}
	}
	if counter == len(p.Payload) {
		return false
	}
	return true
}

func sendPacket(serialPort serial.Port, p Xblock) {

	var sp []byte
	sp = append(sp, p.Proto[0])
	sp = append(sp, byte(p.PacketNum)+1)
	sp = append(sp, byte(byte(255)-byte(p.PacketNum)-1))
	for _, b := range p.Payload {
		sp = append(sp, b)
	}
	sp = append(sp, byte(p.Checksum))
	serialPort.Write(sp)
}
