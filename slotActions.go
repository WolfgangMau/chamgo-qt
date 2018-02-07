package main

import (
	"bytes"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var temp2 []string

var myTime time.Time
var GetUsbListTicker *time.Ticker

func buttonClicked(btn int) {

	switch ActionButtons[btn] {

	case "Select All":
		selectAllSlots(true)
		if populated {
		}

	case "Select None":
		selectAllSlots(false)
		if populated {
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
		log.Printf("clicked on Button: %s\n", ActionButtons[btn])
	}
}

func slotChecked(slot, state int) {
	//log.Printf(" Checked %d - state: %d\n", slot, state)
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
	for i, s := range Slots {
		if s.slot.IsChecked() {
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
			sendSerialCmd(Commands.buttonl + "=" + s.btnl.CurrentText())
		}
	}
	populateSlots()
}

func countSelected() int {
	c := 0
	for _, s := range Slots {
		if s.slot.IsChecked() {
			c++
		}
	}
	return c
}

func clearSlot() {
	c1 := 0
	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			c1++
			log.Printf("clearing %s\n", s.slotl.Text())
			hardwareSlot := i
			if Device == Devices.name[1] {
				hardwareSlot = i + 1
			}
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))
			sendSerialCmd(DeviceActions.clearSlot)
		}
	}
	time.Sleep(time.Millisecond * 50)
	populateSlots()
}

func refreshSlot() {
	populateSlots()
}

func activateSlots() {
	if countSelected() > 1 {
		widgets.QMessageBox_Information(nil, "OK", "please select only one Slot",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			hardwareSlot := i
			if Device == Devices.name[1] {
				hardwareSlot = i + 1
			}
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))
		}
	}
}

//ToDO: implemetation
func mfkey32Slots() {
	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			log.Printf("I should probably calc keys for Slot %d\n", i)
		}
	}
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
	if filename == "" {
		log.Println("no file selected")
		return false
	}

	for i, s := range Slots {
		if s.slot.IsChecked() {
			hardwareSlot := i
			if Device == Devices.name[1] {
				hardwareSlot = i + 1
			}
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))
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

			var p []xblock
			var p1 xblock
			for _, d := range data {
				p1.payload = append(p1.payload, d)

				if len(p1.payload) == 128 {
					p1.proto = []byte{SOH}
					p1.packetNum = len(p)
					p1.packetInv = 255 - p1.packetNum
					p1.checksumm = int(checksum(p1.payload, 0))
					p = append(p, p1)
					p1.payload = []byte("")
				}
			}

			//set chameleon into receiver-mode
			sendSerialCmd(DeviceActions.startUpload)
			if SerialResponse.Code == 110 {
				//start uploading packets
				xmodemSend(p)
			}
		}
	}
	refreshSlot()
	return true
}

func sendPacket(p xblock) {

	var sp []byte
	sp = append(sp, p.proto[0])
	sp = append(sp, byte(p.packetNum)+1)
	sp = append(sp, byte(byte(255)-byte(p.packetNum)-1))
	for _, b := range p.payload {
		sp = append(sp, b)
	}
	sp = append(sp, byte(p.checksumm))
	serialPort.Write(sp)
}

func checksum(b []byte, cs byte) byte {
	for _, d := range b {
		cs = cs + d
	}
	return cs
}

//returns false if all payload-bytes are set to 0xff
func (p xblock) checkPaylod() bool {
	var counter = 0
	for _, b := range p.payload {
		if b == 0xff {
			counter++
		}
	}
	if counter == len(p.payload) {
		return false
	}
	return true
}

func downloadSlots() {
	if countSelected() > 1 {
		widgets.QMessageBox_Information(nil, "OK", "please select only one Slot",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	var data bytes.Buffer
	var (
		success int
		failed int
	)
	var filename string
	for i, s := range Slots {
		hardwareSlot := i
		if Device == Devices.name[1] {
			hardwareSlot = i + 1
		}
		sel := s.slot.IsChecked()
		if sel {
			fileSelect := widgets.NewQFileDialog(nil, 0)
			filename = fileSelect.GetSaveFileName(nil, "save Data from "+s.slotl.Text()+" to File", "", "", "", fileSelect.Options())
			if filename == "" {
				log.Println("no file seleted")
				return
			}
			log.Printf("download a dump from Slot %d into file %s\n", i, filename)
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))

			//set chameleon into receiver-mode
			sendSerialCmd(DeviceActions.startDownload)
			if SerialResponse.Code == 110 {
			success, failed, data  = xmodemRead()

				log.Printf("Success: %d - failed: %d\n", success, failed)
			}
			if _, err := serialPort.Write([]byte{CAN}); err != nil {
				log.Println(err)
				break
			}

			if data.Len() > 0 {
				log.Printf("got %d bytes to write to %s... ", data.Len(), filename)
				// Write file
				fOut, err := os.Create(filename)
				if err != nil {
					log.Println(filename, " - write failed")
					log.Fatalln(err)
				}
				fOut.Write(data.Bytes())
				fOut.Close()

				log.Println(filename, " - write successful")
			} else {
				log.Printf("got only %d bytes - file not written", data.Len())
			}
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

	c := 0
	all := countSelected()

	hardwareSlot := 0
	myProgressBar.zero()
	myProgressBar.widget.SetRange(c, all)
	for sn, s := range Slots {
		//update single slot
		if s.slot.IsChecked() {
			c++
			myProgressBar.update(c)
			if Device == Devices.name[1] {
				hardwareSlot = sn + 1
			} else {
				hardwareSlot = sn
			}

			log.Printf("read data for Slot %d\n", sn+1)
			sendSerialCmd(DeviceActions.selectSlot + strconv.Itoa(hardwareSlot))
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

			sendSerialCmd(DeviceActions.getButtonl)
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

func checkForDevices() {
	GetUsbListTicker = time.NewTicker(time.Millisecond * 5000)
	go func() {
		for myTime = range GetUsbListTicker.C {
			if !Connected {
				serialPorts, err := getSerialPorts()
				if err != nil {
					log.Println(err)
				}
				if len(serialPorts) > 0 && serialPortSelect.CurrentText() != serialPorts[SelectedPortId] {
					serialPortSelect.Clear()
					serialPortSelect.AddItems(serialPorts)
					serialPortSelect.SetCurrentIndex(SelectedPortId)
					serialPortSelect.Repaint()

					deviceSelect.SetCurrentIndex(SelectedDeviceId)
					deviceSelect.Repaint()
					//Device = deviceSelect.CurrentText()
					//Commands.load(Device)
					//DeviceActions.load(Device)
				} else {
					if len(serialPorts) == 0 {
						serialPortSelect.Clear()
					}
				}
			} else {
				GetUsbListTicker.Stop()
			}
		}
	}()
}

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

func (pb *progressBar) update(c int) {
	pb.widget.SetValue(c)
}

func (pb *progressBar) zero() {
	pb.widget.Reset()
}
