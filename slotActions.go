package main

import (
	"github.com/therecipe/qt/widgets"
	//"github.com/chrizzzzz/go-xmodem/xmodem"
	"log"
	"strconv"
	"strings"
	"os"
	"io/ioutil"
	"time"
)

var temp2 []string

/*var temp string
var myTime time.Time
var GetSlotTicker *time.Ticker*/

func buttonClicked(btn int) {

	switch ActionButtons[btn] {

	case "Select All":
		selectAllSlots(true)
		if populated {
			//GetSlotTicker.Stop()
		}

	case "Select None":
		selectAllSlots(false)
		if populated {
			//GetSlotTicker.Stop()
		}

	case "Apply":
		applySlot()

	case "Clear":
		clearSlot()

	case "Refresh":
		refreshSlot()

	case "Set Active":
		activateSlots()

	case "mfkey32":
		mfkey32Slots()

	case "Upload":
		uploadSlots()

	case "Download":
		downloadSlots()

	default:
		log.Printf("clicked on Button %s\n", ActionButtons[btn])
	}
}

func slotChecked(slot, state int) {
	log.Printf(" Checked %d - state: %d\n", slot, state)
	if state == 2 && Connected {
		if Device == Devices.name[1] {
			//RevG's first Slot is 1 and Last Slot is 8
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(slot+1))
		} else {
			//RevE's first Slot is 0 and Last Slot is 7
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(slot))
		}
	}
	Slots[slot].slot.Repaint()
}

func selectAllSlots(b bool) {
	for _, s := range Slots {
		s.slot.SetChecked(b)
		s.slot.Repaint()
	}
}

func applySlot() {
	//GetSlotTicker.Stop()
	for i, s := range Slots {
		if s.slot.IsChecked() {
			log.Printf("********************\nupdating %s\n", s.slotl.Text())
			hardwareSlot := i
			if Device == Devices.name[1] {
				hardwareSlot = i + 1
			}
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))
			//select slot
			sendSerialCmd(Commands.config + "=" + s.mode.CurrentText())
			//set mode
			sendSerialCmd(Commands.config + "=" + s.mode.CurrentText())
			//set uid
			sendSerialCmd(Commands.uid + "=" + s.uid.Text())
			//set  button short
			sendSerialCmd(Commands.button + "=" + s.btns.CurrentText())
			//set button long
			sendSerialCmd(Commands.lbutton + "=" + s.btnl.CurrentText())
		}
	}
	populateSlots()
}

func countSelected() int {
	c:=0
	for _, s := range Slots {
		if s.slot.IsChecked() {
			c++
		}
	}
	return c
}

func clearSlot() {
	c1:=0
	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			c1++
			log.Printf("********************\nclearing %s\n", s.slotl.Text())
			hardwareSlot := i
			if Device == Devices.name[1] {
				hardwareSlot = i + 1
			}
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))
			sendSerialCmd(DeviceActions.clearSlot)
		}
	}
	populateSlots()
}

func refreshSlot() {
	populateSlots()
}

func activateSlots() {
	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			log.Printf("I should probably activate Slot %d\n", i)
		}
	}
}

func mfkey32Slots() {
	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			log.Printf("I should probably calc keys for Slot %d\n", i)
		}
	}
}

type packet struct {
	proto byte
	block int
	rblocks int
	data  []byte
	chk   byte
}


func uploadSlots() {
	var filename string
	fileSelect := widgets.NewQFileDialog(nil, 0)
	filename = fileSelect.GetOpenFileName(nil, "open Dump", "", "", "", fileSelect.Options())
    if filename == ""{
    	log.Println("no file selöeted")
    	return
	}

	for i, s := range Slots {
		if s.slot.IsChecked() {
			log.Printf("********************\nupdating %s\n", s.slotl.Text())
			hardwareSlot := i
			if Device == Devices.name[1] {
				hardwareSlot = i + 1
			}
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))
			log.Printf("I should probably upoload %s to Slot %d\n", filename, i)
			/**
							XMODEM 128 byte blocks
							----------------------
				SENDER                                      RECEIVER
														<-- NAK
				SOH 01 FE Data[128] CSUM                -->
														<-- ACK
				SOH 02 FD Data[128] CSUM                -->
														<-- ACK
				SOH 03 FC Data[128] CSUM                -->
														<-- ACK
				SOH 04 FB Data[128] CSUM                -->
														<-- ACK
				SOH 05 FA Data[100] CPMEOF[28] CSUM     -->
														<-- ACK
				EOT                                     -->
														<-- ACK

			 */

			 // very basic implementation of a xmodem-sender of mine
			 // buggy on windows32 (uploads, but no feedback - app freezes
			 // because the chameleon is not stopping from xmodem-receiver-mode)
			 // the usb-device must be removed to cut connection
			 // but works mostly on osx an linux ...
			 // the lack of responses is the most problesm, since I was able to get
			 // a bi-directional data-exchange working (I don't get a NAK, ACK from receiver)
			 // so just a 'fire and forget' mission for the data, but mostly it works
			 // ToDo: needs to be tested with other serial libs - reader is needed for download also!
			 // ToDo: in order to get closer to the xmodem specification:
			 // 	- fill eventually 'not filled' data-blocks with EOF's
			 // 	- retry (10 times) on transmision-failure -> (blocked by bi-direction-issue)
			 //     - make a library out of it

			// Open file
			log.Printf("loading file %s\n", filename)
			fIn, err := os.Open(filename)
			if err != nil {
				log.Fatalln(err)
			}
			//readfile into buffer
			data, err := ioutil.ReadAll(fIn)
			if err != nil {
				log.Println(err)
			}
			fIn.Close()

			var p []packet
			var p1 packet
			//build 128byte packages (works at this stage only for 1k and 4k)
			for _,d := range data {
				//temp := byte(d)
					//temp := byte(d)
					p1.data = append(p1.data, d)

					if len(p1.data) == 128 {
						p1.proto = 0x01
						p1.block = len(p)
						p1.rblocks = 255-len(p)
						p1.chk = checksum(p1.data, 0)
						p = append(p, p1)
						p1.data=[]byte("")
					}
			}

			//set chameleon into receiver-mode
			sendSerialCmd(DeviceActions.startUpload)

			//re-establish a fresh connection
			err = serialPort.Close()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Millisecond * 200)
			err = connectSerial(SerialDevice1)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Millisecond * 500)
			// send NAK byte
			//var nak []byte
			//nak = append(nak,0x06)
			//_,err = serialPort.Write(nak)


			// send ACK byte
			//var ack []byte
			//ack = append(ack,0x15)
			//_,err = serialPort.Write(ack)

			//send all packets
			//all:=len(p)
			//myProgressBar.widget.SetStatusTip("uploadiong file "+filename)
			//myProgressBar.widget.SetRange(0,all)
			for _,sp := range p {
				sendPacket(sp)
				time.Sleep(time.Millisecond * 25)
				//myProgressBar.update(i)
			}
		    log.Println("upload done")

			//send EOF byte
			//var eof []byte
			//eof = append(eof,0x1a)
			//_,err = serialPort.Write(eof)
			//if err != nil {
			//	log.Println(err)
			//}


		    //send EOT byte
			var eot []byte
			eot = append(eot,0x04)
			serialPort.Write(eot)

			time.Sleep(time.Millisecond * 100)
			serialPort.Write(eot)

			if err != nil {
				log.Println(err)
			}

			//send CAN byte
			//var can []byte
			//can = append(can,0x18)
			//_,err = serialPort.Write(can)

			//for i2:=0; i2< 10; i2++ {
			//	for i:=0; i<131; i++ {
			//		serialPort.Write(can)
			//	}
			//}

			//clear up serial buffers ...
			//serialPort.ResetOutputBuffer()
			//serialPort.ResetInputBuffer()

			//close serial
			err = serialPort.Close()
			if err != nil {
				log.Println(err)
			}
			//start serial and cut xmodem-receiver (again 4 windows)
			//myProgressBar.update(all)

			err = connectSerial(SerialDevice1)
			if err != nil {
				log.Println(err)
			}

			populateSlots()
		}
	}
}

func sendPacket(p packet){

	var sp []byte
	sp = append(sp, p.proto)
	sp = append(sp, byte(p.block)+1)
	sp = append(sp, byte(byte(255)-byte(p.block)-1))
	for _,b := range p.data {
		sp = append(sp, b)
	}
	sp = append(sp, p.chk)
	serialPort.Write(sp)

	//time.Sleep(time.Millisecond * 10)
	var resp []byte
	//time.Sleep(time.Millisecond * 100)
	i, err := serialPort.Read(resp)
	if err != nil {
		log.Println(err)
	}
	if len(resp)>0 {
		log.Printf("got response from receiver: %x (%d)\n", resp,i)
	}

}

func  checksum(b []byte, cs byte) byte {
	for _,d := range b {
		cs = cs + d
	}
	return cs
}

func downloadSlots() {
	var filename string
	fileSelect := widgets.NewQFileDialog(nil, 0)
	filename = fileSelect.GetOpenFileName(nil, "save Dump", "", "", "", fileSelect.Options())
	if filename == ""{
		log.Println("no file selöeted")
		return
	}

	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			log.Printf("should probably download a dump from Slot %d into file %s\n", i, filename)
		}
	}

}

func populateSlots() {
	if !Connected {
		return
	}

	if populated == false {
		//ToDo: error-handling
		sendSerialCmd(DeviceActions.getModes)
		TagModes = strings.Split(SerialResponse.Payload, ",")
		//ToDo: error-handling
		sendSerialCmd(DeviceActions.getButtons)
		TagButtons = strings.Split(SerialResponse.Payload, ",")
		populated = true
	}

	c:=0
	all:=countSelected()

	softSlot:=0
	myProgressBar.zero()
	myProgressBar.widget.SetRange(c,all)
	for sn, s := range Slots {
		//update single slot
		if s.slot.IsChecked() {
			c++
			myProgressBar.update(c)
			log.Printf("update %d\n",c)
			if Device == Devices.name[1] {
				softSlot = sn + 1
			} else {
				softSlot = sn
			}

			sendSerialCmd(DeviceActions.selectSlot+strconv.Itoa(softSlot))
			//get slot uid
			sendSerialCmd(DeviceActions.getUid)
			uid := SerialResponse.Payload
			//set uid to lineedit
			s.uid.SetText(uid)

			sendSerialCmd(DeviceActions.getSize)
			size := SerialResponse.Payload

			s.size.SetText(size)

			sendSerialCmd(DeviceActions.getMode)
			mode := SerialResponse.Payload
			_, modeindex := getPosFromList(mode, TagModes)
			s.mode.Clear()
			s.mode.AddItems(TagModes)
			s.mode.SetCurrentIndex(modeindex)
			s.mode.Repaint()

			sendSerialCmd(DeviceActions.getButton)
			buttonl := SerialResponse.Payload
			_, buttonlindex := getPosFromList(buttonl, TagButtons)
			s.btnl.Clear()
			s.btnl.AddItems(TagButtons)
			s.btnl.SetCurrentIndex(buttonlindex)
			s.btnl.Repaint()

			// ToDo: currently mostly faked - currently not implemented in my revG
			//unlear about RButton & LButton short and long -> 4 scenarios?
			sendSerialCmd(DeviceActions.getButton)
			buttons := SerialResponse.Payload
			_, buttonsindex := getPosFromList(buttons, TagButtons)
			s.btns.Clear()
			s.btns.AddItems(TagButtons)
			s.btns.SetCurrentIndex(buttonsindex)
			s.btns.Repaint()
		}
	}
}
//
//func checkCurrentSelection() {
//	GetSlotTicker = time.NewTicker(time.Millisecond * 2000)
//	var softSlot int
//	go func() {
//		for myTime = range GetSlotTicker.C {
//			sendSerialCmd(DeviceActions.selectedSlot)
//			selected := SerialResponse.Payload
//			if Device == Devices.name[1] {
//				hardSlot, _ := strconv.Atoi(selected)
//				softSlot = hardSlot - 1
//			} else {
//				hardSlot, _ := strconv.Atoi(strings.Replace(selected, "NO.", "", 1))
//				softSlot = hardSlot
//			}
//			log.Printf("Tick at %s - Current Selected Slot: %d\n\n", myTime, softSlot+1)
//			for i, s := range Slots {
//				if s.slot.IsChecked() && i != softSlot {
//					s.slot.SetChecked(false)
//				} else {
//					if !s.slot.IsChecked() && i == softSlot && populated {
//						s.slot.SetChecked(true)
//					}
//				}
//			}
//		}
//	}()
//}

func getPosFromList(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}

	return
}

func (pb progressBar) update(c int)  {
	pb.widget.SetValue(c)
}

func (pb *progressBar) zero() {
	pb.widget.Reset()
	//pb.widget.SetValue(0)
	//pb.widget.Repaint()
}