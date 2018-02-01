package main

import (
	"github.com/therecipe/qt/widgets"
	//"github.com/chrizzzzz/go-xmodem/xmodem"
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
		log.Printf("Clicked on Button %s\n", ActionButtons[btn])
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
}

func clearSlot() {
	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			log.Printf("I should probably clear settings to Slot %d\n", i)
		}
	}
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

	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
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

			>>> csum = modem.calc_checksum('hello')
			>>> csum = modem.calc_checksum('world', csum)
			>>> hex(csum)
			'0x3c'

			package main
			import (
				"fmt"
			)

			func main() {
				var my []byte
				my = []byte("helloworld")
				c := byte(0)
				for _,b := range my {
					c = c+b
				}
				fmt.Printf("chk: 0x%x\n", c)
			}
			 */
			////send file
			//
			//
			//// Open file
			log.Printf("loading file %s\n", filename)
			fIn, err := os.Open(filename)
			if err != nil {
				log.Fatalln(err)
			}

			data, err := ioutil.ReadAll(fIn)
			if err != nil {
				log.Println(err)
			}
			fIn.Close()

			var p []packet
			var p1 packet
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
			log.Printf("packets : %d  -  dataLen: %d\n", len(p), len(p[0].data) )
			//
			//sendSerialCmd(DeviceActions.startUpload)
			//err = serialPort.Close()
			//if err != nil {
			//	log.Println(err)
			//}
			//log.Println("port closed")
			//err = connectSerial(SerialDevice1)
			//if err != nil {
			//	log.Println(err)
			//}
			//log.Println("port opend")
			//
			//log.Println(filename, "start upload")
			//// Send file
			//var pl []byte
			//xmodem.SendBlock(serialPort,0, pl, 1)
			//err = xmodem.ModemSend(serialPort, data)
			//log.Println(filename, "finished upload")
			//if err != nil {
			//	log.Println(err)
			//}
			//err := serialPort.Close()
			//if err != nil {
			//	log.Println(err)
			//}
			//log.Println(filename, "sent successful")
			//ex, err := os.Executable()
			//if err != nil {
			//	panic(err)
			//}
			//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			//time.Sleep(time.Second * 1)
			//cmd := exec.Command(dir+string(os.PathSeparator)+"xmutil", SerialDevice1, filename)
			//var out bytes.Buffer
			//cmd.Stdout = &out
			//err = cmd.Run()
			//if err != nil {
			//	log.Fatal(err)
			//}
			//time.Sleep(time.Millisecond * 100)
			//err = connectSerial(SerialDevice1)
			//if err != nil {
			//	log.Println(err)
			//}
		}
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

	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			log.Printf("I should probably download a dump from Slot %d into file %s\n", i, filename)
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

	softSlot:=0
	for sn, s := range Slots {
		//update single slot
		if s.slot.IsChecked() {
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

