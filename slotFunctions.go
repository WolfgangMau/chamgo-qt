package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/WolfgangMau/chamgo-qt/nonces"
	"github.com/WolfgangMau/chamgo-qt/xmodem"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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
	log.Printf(" Checked %d - state: %d\n", slot, state)
	if state == 2 && Connected {
		sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(slot+Cfg.Device[SelectedDeviceId].Config.Slot.Offset))
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
			hardwareSlot := i + Cfg.Device[SelectedDeviceId].Config.Slot.Offset
			sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(hardwareSlot))
			//select slot
			sendSerialCmd(Cfg.Device[SelectedDeviceId].CmdSet["config"] + "=" + s.mode.CurrentText())
			//set mode
			sendSerialCmd(Cfg.Device[SelectedDeviceId].CmdSet["config"] + "=" + s.mode.CurrentText())
			//set uid
			sendSerialCmd(Cfg.Device[SelectedDeviceId].CmdSet["uid"] + "=" + s.uid.Text())
			//set  button short
			sendSerialCmd(Cfg.Device[SelectedDeviceId].CmdSet["button"] + "=" + s.btns.CurrentText())
			//set button long
			sendSerialCmd(Cfg.Device[SelectedDeviceId].CmdSet["buttonl"] + "=" + s.btnl.CurrentText())
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
			hardwareSlot := i + Cfg.Device[SelectedDeviceId].Config.Slot.Offset
			sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(hardwareSlot))
			sendSerialCmd(DeviceActions.ClearSlot)
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
			hardwareSlot := i + Cfg.Device[SelectedDeviceId].Config.Slot.Offset
			sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(hardwareSlot))
			Cfg.Device[SelectedDeviceId].Config.Slot.Selected = hardwareSlot
		}
	}
}

//ToDO: implemetation
func mfkey32Slots() {
	if !Connected || countSelected() < 1 {
		if !Connected {
			return
		}
		widgets.QMessageBox_Information(nil, "OK", "please select at least one Slot\nwhich was set to DETECTION",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	detectionCmd, ok := Cfg.Device[SelectedDeviceId].CmdSet["detection"]
	if !ok {
		widgets.QMessageBox_Information(nil, "OK", "Sorry, but this Device hs not set a 'detection' cmd!",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	for i, s := range Slots {
		sel := s.slot.IsChecked()
		if sel {
			hardwareSlot := i + Cfg.Device[SelectedDeviceId].Config.Slot.Offset
			sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(hardwareSlot))
			serialMonitor.AppendPlainText("-> " + strings.Replace(strings.Replace(detectionCmd, "\r", "", -1), "\n", "", -1))

			// send cmd ang get the expected 218 bytes (208 nonce + 2 crc + 8 cmd-response (100:OK\n\r)
			SerialSendOnly(detectionCmd + "?")
			buff := GetSpecificBytes(218)
			//buffer should be empty - only to get sure
			SerialPort.ResetInputBuffer()

			//responsecode := strings.Replace(strings.Replace(string(buff[len(buff)-8:]), "\r", "", -1), "\n", "", -1)
			log.Println("len enc: ", len(buff))
			buff = nonces.DecryptData(buff[0:len(buff)-10], 123321, 208)
			uid := buff[0:4]
			serialMonitor.AppendPlainText(fmt.Sprintf("uid: %x\n", uid))

			noncemap := nonces.ExtractNonces(buff)
			var skey string
			if len(noncemap) > 0 {
				MyTabs.SetCurrentIndex(1)
				serialMonitor.AppendPlainText(fmt.Sprintf("Fond %d nonces for UID: %04X - test possible comboinations ...", len(noncemap), uid))
				log.Println("  UID      NT0      NR0      AR0      NT1      NR1      AR1")
				for i1 := 0; i1 < len(noncemap); i1++ {
					for i2 := 0; i2 < len(noncemap); i2++ {
						if i1 == i2 || i1 > i2 {
							continue
						} else {
							if noncemap[i1].Key == noncemap[i2].Key {
								if noncemap[i1].Key == 0x60 {
									skey = "A"
								} else {
									skey = "B"
								}
								args := []string{hex.EncodeToString(uid), hex.EncodeToString(noncemap[i1].Nt), hex.EncodeToString(noncemap[i1].Nr), hex.EncodeToString(noncemap[i1].Ar), hex.EncodeToString(noncemap[i2].Nt), hex.EncodeToString(noncemap[i2].Nr), hex.EncodeToString(noncemap[i2].Ar)}
								log.Printf("%04X %04X %04X %04X %04X %04X %04X\n", uid, noncemap[i1].Nt, noncemap[i1].Nr, noncemap[i1].Ar, noncemap[i2].Nt, noncemap[i2].Nr, noncemap[i2].Ar)
								res, err := execCmd("mfkey32v2", args)
								if err != nil {
									log.Println(err)
								} else {
									if strings.Contains(res, "Found Key") {
										key := strings.Split(res, "[")[1]
										key = key[:12]
										serialMonitor.AppendPlainText(fmt.Sprintf("Slot %d: Possible Key %s for Nonces on  Blocks %d & %d = %s", i+1, skey, noncemap[i1].Sector, noncemap[i2].Sector, key))

									}
								}
							}
						}
					}
				}
			}
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
			hardwareSlot := i + Cfg.Device[SelectedDeviceId].Config.Slot.Offset
			sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(hardwareSlot))
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

			var p []xmodem.Xblock
			var p1 xmodem.Xblock

			for _, d := range data {
				p1.Payload = append(p1.Payload, d)

				if len(p1.Payload) == 128 {
					p1.Proto = []byte{xmodem.SOH}
					p1.PacketNum = len(p)
					p1.PacketInv = 255 - p1.PacketNum
					p1.Checksum = int(xmodem.Checksum(p1.Payload, 0))
					p = append(p, p1)
					p1.Payload = []byte("")
				}
			}

			//set chameleon into receiver-mode
			sendSerialCmd(DeviceActions.StartUpload)
			if SerialResponse.Code == 110 {
				//start uploading packets
				xmodem.Send(SerialPort, p)
			}
		}
	}
	refreshSlot()
	return true
}

func downloadSlots() {
	if countSelected() > 1 {
		widgets.QMessageBox_Information(nil, "OK", "please select only one Slot",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}

	var (
		success int
		failed  int
	)
	var filename string
	var data bytes.Buffer

	for i, s := range Slots {
		hardwareSlot := i + Cfg.Device[SelectedDeviceId].Config.Slot.Offset
		sel := s.slot.IsChecked()
		if sel {
			fileSelect := widgets.NewQFileDialog(nil, 0)
			filename = fileSelect.GetSaveFileName(nil, "save Data from "+s.slotl.Text()+" to File", "", "", "", fileSelect.Options())
			if filename == "" {
				log.Println("no file seleted")
				return
			}
			log.Printf("download a dump from Slot %d into file %s\n", i, filename)
			sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(hardwareSlot))

			//set chameleon into receiver-mode
			sendSerialCmd(DeviceActions.StartDownload)
			if SerialResponse.Code == 110 {
				success, failed, data = xmodem.Receive(SerialPort)

				log.Printf("Success: %d - failed: %d\n", success, failed)
			}
			if _, err := SerialPort.Write([]byte{xmodem.CAN}); err != nil {
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
		sendSerialCmd(DeviceActions.GetModes)
		TagModes = strings.Split(SerialResponse.Payload, ",")
		//ToDo: error-handling
		sendSerialCmd(DeviceActions.GetButtons)
		TagButtons = strings.Split(SerialResponse.Payload, ",")
		populated = true
	}

	c := 0

	hardwareSlot := 0
	for sn, s := range Slots {
		//update single slot
		if s.slot.IsChecked() {
			c++
			hardwareSlot = sn + Cfg.Device[SelectedDeviceId].Config.Slot.Offset

			log.Printf("read data for Slot %d\n", sn+1)
			sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(hardwareSlot))
			//get slot uid
			sendSerialCmd(DeviceActions.GetUid)
			uid := SerialResponse.Payload
			//set uid to lineedit
			s.uid.SetText(uid)

			sendSerialCmd(DeviceActions.GetSize)
			size := SerialResponse.Payload

			s.size.SetText(size)

			sendSerialCmd(DeviceActions.GetMode)
			mode := SerialResponse.Payload
			_, modeindex := getPosFromList(mode, TagModes)
			s.mode.Clear()
			s.mode.AddItems(TagModes)
			s.mode.SetCurrentIndex(modeindex)
			s.mode.Repaint()

			sendSerialCmd(DeviceActions.GetButtonl)
			buttonl := SerialResponse.Payload
			_, buttonlindex := getPosFromList(buttonl, TagButtons)
			s.btnl.Clear()
			s.btnl.AddItems(TagButtons)
			s.btnl.SetCurrentIndex(buttonlindex)
			s.btnl.Repaint()

			// ToDo: currently mostly faked - currently not implemented in my revG
			// unlear about RButton & LButton short and long -> 4 scenarios?
			// but works on RevG
			sendSerialCmd(DeviceActions.GetButton)
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

func execCmd(cmdstr string, args []string) (res string, err error) {
	res = ""
	err = nil

	//set local path
	os.Setenv("PATH", os.Getenv("PATH")+":"+os.Getenv("PWD")+"/bin/")

	// Create an *exec.Cmd
	cmd := exec.Command(cmdstr, args...)

	// Stdout buffer
	cmdOutput := &bytes.Buffer{}
	// Attach buffer to command
	cmd.Stdout = cmdOutput

	// Execute command
	//log.Printf("run Cmd: %s %s\n",cmd.Path,cmd.Args)
	err = cmd.Run()
	if err != nil {
		log.Printf("Error: %s\n", err)
		return res, err
	}

	// Only output the commands stdout
	res = string(cmdOutput.Bytes())
	//log.Printf("RES: %s\n",res)
	return res, err
}
