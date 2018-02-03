package main

import (
	"github.com/therecipe/qt/widgets"
	"log"
	"strconv"
	"strings"
	"os"
	"io/ioutil"
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


func uploadSlots() bool {
	if countSelected() > 1 {
		widgets.QMessageBox_Information(nil, "OK", "please select only one Slot",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return false
	}
	var filename string
	fileSelect := widgets.NewQFileDialog(nil, 0)
	filename = fileSelect.GetOpenFileName(nil, "open Dump", "", "", "", fileSelect.Options())
    if filename == ""{
    	log.Println("no file selöeted")
    	return false
	}

	for i, s := range Slots {
		if s.slot.IsChecked() {
			log.Printf("********************\nupdating %s\n", s.slotl.Text())
			hardwareSlot := i
			if Device == Devices.name[1] {
				hardwareSlot = i + 1
			}
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))
			log.Printf("upoload %s to Slot %d\n", filename, i)
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
			oBuffer := make([]byte, 1)
			for _, d := range data {
				p1.data = append(p1.data, d)

				if len(p1.data) == 128 {
					p1.proto = 0x01
					p1.block = len(p)
					p1.rblocks = 255 - len(p)
					p1.chk = checksum(p1.data, 0)
					p = append(p, p1)
					p1.data = []byte("")
				}
			}

			//set chameleon into receiver-mode
			sendSerialCmd(DeviceActions.startUpload)
			if SerialResponse.Code == 110 {

				//send NAK byte / init transfer
				//var nak []byte
				//nak = append(nak, 0x04)
				//serialPort.Write(nak)
				//if _,err = serialPort.Read(oBuffer); err != nil {
				//	log.Println(err)
				//}
				//if oBuffer[0] != 0x06 {
				//	log.Printf("nexpectedanswer to NAK: 0x%X\n", oBuffer[0])
				//}

				//start uploading packets
				failure := 0
				success := 0
				for _, sp := range p {
					var reSend bool = true
					for reSend {
						sendPacket(sp)
						if _,err = serialPort.Read(oBuffer); err != nil {
							log.Println(err)
						} else {
							switch oBuffer[0] {
								case 0x015: // NAK
									log.Printf("resend packet %d\n",sp.block)
									reSend = true
									failure++
								case 0x06: // ACK
									reSend = false
									success++
								default:
									log.Printf("unexspected answer(0x%X) for packet %d\n",oBuffer[0],sp.block)
									reSend = false
							}
						}
					}
					//myProgressBar.update(i)
				}
				log.Printf("upload done - Success: %d - Failures: %d\n", success, failure)

				//send EOT byte
				var eot []byte
				eot = append(eot, 0x04)
				serialPort.Write(eot)
				if _,err = serialPort.Read(oBuffer); err != nil {
					log.Println(err)
				}
				if oBuffer[0] != 0x06 {
					log.Printf("nexpectedanswer to EOT: 0x%X\n", oBuffer[0])
				}

				////send CAN byte
				//var can []byte
				//can = append(can,0x18)
				//_,err = serialPort.Write(can)
				//if _,err = serialPort.Read(oBuffer); err != nil {
				//	log.Println(err)
				//}
				//if oBuffer[0] != 0x06 {
				//	log.Printf("unexpected answer to CAN: 0x%X\n",oBuffer[0])
				//}
			}
		}
	}
	refreshSlot()
	return true
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