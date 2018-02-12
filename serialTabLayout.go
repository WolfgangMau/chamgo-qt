package main

import (
	"bufio"
	"fmt"
	"github.com/therecipe/qt/widgets"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	serialSendButton *widgets.QPushButton
	serialMonitor    *widgets.QPlainTextEdit
	serialPortSelect *widgets.QComboBox
	deviceSelect     *widgets.QComboBox
)

func serialTab() *widgets.QWidget {
	serialTabLayout := widgets.NewQHBoxLayout()
	leftTabLayout := widgets.NewQVBoxLayout()
	serialTabPage := widgets.NewQWidget(nil, 0)

	/********************************************** Serial Connect *********************************************/

	serialPorts, _ := getSerialPorts()

	serConLayout := widgets.NewQFormLayout(nil)

	deviceSelect = widgets.NewQComboBox(nil)
	deviceSelect.AddItems(getDeviceNames())
	deviceSelect.SetCurrentIndex(SelectedDeviceId)
	deviceSelect.SetFixedWidth(160)

	serialPortSelect = widgets.NewQComboBox(nil)
	serialPortSelect.AddItems(serialPorts)
	serialPortSelect.SetFixedWidth(160)
	serialPortSelect.SetCurrentIndex(SelectedPortId)

	serialConnectButton := widgets.NewQPushButton2("Connect", nil)

	serialDeviceInfo := widgets.NewQLabel(nil, 0)
	serialDeviceInfo.SetText("not Connected")

	serConLayout.AddWidget(deviceSelect)
	serConLayout.AddWidget(serialPortSelect)
	serConLayout.AddWidget(serialConnectButton)
	serConLayout.AddWidget(serialDeviceInfo)

	serialConnectGroup := widgets.NewQGroupBox2("Serial Connection", nil)
	serialConnectGroup.SetLayout(serConLayout)
	serialConnectGroup.SetFixedSize2(220, 180)

	macroGroupLayout := widgets.NewQHBoxLayout()
	macroGroup := widgets.NewQGroupBox2("Command Macros", nil)
	macroGroup.SetFixedWidth(220)
	macroSelect := widgets.NewQComboBox(macroGroup)
	macroSelect.AddItems(getFilesInFolder("macros", ".cmds"))
	macroGroupLayout.AddWidget(macroSelect, 1, 0x0020)
	macroSend := widgets.NewQPushButton2("execute", nil)
	macroGroupLayout.AddWidget(macroSend, 1, 0x0020)

	macroGroup.SetLayout(macroGroupLayout)

	leftTabLayout.AddWidget(serialConnectGroup, 1, 0x0020)
	leftTabLayout.AddWidget(macroGroup, 1, 0x0020)
	serialTabLayout.AddLayout(leftTabLayout, 0)

	serialConnectButton.ConnectClicked(func(checked bool) {
		Commands := Cfg.Device[SelectedDeviceId].CmdSet
		//for n,c := range Commands {
		//	log.Printf("Command Name: %s -> %s\n",n,c)
		//}
		//Commands.load(deviceSelect.CurrentText())

		if serialConnectButton.Text() == "Connect" {

			err := connectSerial(SerialDevice1)
			if err != nil {
				widgets.QMessageBox_Information(nil, "OK", "can't connect to Serial\n"+string(err.Error()),
					widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				log.Println("error on connect: ", err)
			} else {
				dn := Cfg.Device[SelectedDeviceId].Name
				DeviceActions.load(dn)
				if len(DeviceActions.getUid) <= 0 {
					log.Println("no action for 'getUid!?' ", DeviceActions.getUid)
				}

				//ask for the device-version
				sendSerialCmd(Commands["version"] + "?")

				if SerialResponse.Code >= 100 {
					serialConnectButton.SetText("Disconnect")
					serialSendButton.SetDisabled(false)
					serialSendButton.Repaint()
					//checkForDevices()
				}
				//web got a expected answer from the VERSION(MY) cmd
				if SerialResponse.Code == 101 {
					serialDeviceInfo.SetText("Connected\n" + deviceInfo(SerialResponse.Payload))
					Connected = true
					Statusbar.ShowMessage("Connected to Port: "+serialPortSelect.CurrentText()+" - Device: "+Cfg.Device[SelectedDeviceId].Cdc+" - Firmware: "+deviceInfo(SerialResponse.Payload), 0)
					if SerialResponse.Code == 101 {
						buttonClicked(0)
						buttonClicked(4)
						buttonClicked(1)
					}
				} else {
					widgets.QMessageBox_Information(nil, "OK", "no Version Response from Device!",
						widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				}
			}
		} else {
			err := serialPort.Close()
			if err == nil {
				Cfg.Save()
				serialConnectButton.SetText("Connect")
				serialSendButton.SetDisabled(true)
				serialSendButton.Repaint()
				serialDeviceInfo.SetText("not Connected")
				Statusbar.ShowMessage("not Connected", 0)
				Connected = false
				serialPort.Close()

			}
		}

	})

	macroSend.ConnectClicked(func(checked bool) {
		log.Printf("execute macro %s\n", macroSelect.CurrentText())
		cmds := readFileLines(macroSelect.CurrentText())
		if len(cmds) > 0 {
			for _, c := range cmds {
				if strings.Contains(strings.ToLower(c), "detectionmy?") {
					serialMonitor.AppendPlainText("<- " + strings.Replace(strings.Replace(c, "\r", "", -1), "\n", "", -1))
					_, err := serialPort.Write([]byte(strings.ToUpper(c) + "\r\n"))
					if err != nil {
						log.Println(err)
					}
					time.Sleep(time.Millisecond * time.Duration(Cfg.Device[SelectedDeviceId].Config.Serial.WaitForReceive))
					n := 1
					c := 0
					buff := make([]byte, 512)
					for c <= 0 {
						n, err = serialPort.Read(buff)
						if err != nil {
							log.Println(err)
						}
						c = c + n
					}
					log.Printf("len enc: %d\n", len(buff[0:c-10]))
					buff2 := DecryptData(buff[0:c-10], 123321, 208)
					uid := buff2[0:4]
					empty := buff2[4:15]
					log.Printf("uid: %x   crc: %x   empty: %x\n", uid, empty[0:1],empty[1:])
					nonces := extractNonces(buff2)
					log.Printf("found %d nonces\n%v\n", len(nonces), nonces)
					responsecode := strings.Replace(strings.Replace(string(buff[c-8:c]), "\r", "", -1), "\n", "", -1)
					serialMonitor.AppendPlainText(fmt.Sprintf("-> %s\nlen: all: %d\nuid: %x\nbuff (%d): %x\n", responsecode, c, uid, len(buff2), buff2))
					serialMonitor.Repaint()
				} else {
					sendSerialCmd(strings.Replace(strings.Replace(c, "\r", "", -1), "\n", "", -1))
					time.Sleep(time.Millisecond * time.Duration(Cfg.Device[SelectedDeviceId].Config.Serial.WaitForReceive))
					serialMonitor.AppendPlainText(fmt.Sprintf("->Code: %d  String: %s Payload: %s\n", SerialResponse.Code, SerialResponse.String, SerialResponse.Payload))
					serialMonitor.Repaint()
				}
			}
		}
	})
	/********************************************** Serial Monitor *********************************************/

	serMonitorLayout := widgets.NewQVBoxLayout()
	serSendLayout := widgets.NewQHBoxLayout()

	serialMonitor = widgets.NewQPlainTextEdit(nil)
	serialMonitor.AppendPlainText("")
	serialMonitor.SetFixedHeight(380)
	serialMonitor.SetReadOnly(true)

	serialSendButton = widgets.NewQPushButton2("send", nil)
	serialSendButton.SetDisabled(true)
	serialSendButton.SetFixedWidth(80)

	serialSendTxt := widgets.NewQLineEdit(nil)
	serialSendTxt.ConnectReturnPressed(serialSendButton.Click)
	serialSendTxt.SetTabOrder(serialSendButton, nil)

	serSendLayout.AddWidget(serialSendTxt, 0, 0)
	serSendLayout.AddWidget(serialSendButton, 1, 0)

	serMonitorLayout.AddWidget(serialMonitor, 0, 0)
	serMonitorLayout.AddLayout(serSendLayout, 0)

	serialSendGroup := widgets.NewQGroupBox2("Serial Terminal", nil)
	serialSendGroup.SetLayout(serMonitorLayout)

	serialTabLayout.AddWidget(serialSendGroup, 0, 0x0020)
	serialTabLayout.SetAlign(33)

	serialSendButton.ConnectClicked(func(checked bool) {
		if serialSendTxt.Text() != "" {
			sendSerialCmd(serialSendTxt.Text())
			if SerialResponse.Code >= 100 {
				serialMonitor.AppendPlainText("-> " + SerialResponse.Cmd)
				serialMonitor.AppendPlainText("<- " + strconv.Itoa(SerialResponse.Code) + " " + SerialResponse.String)
				if SerialResponse.Payload != "" {
					serialMonitor.AppendPlainText("<- " + SerialResponse.Payload)
				}
				serialMonitor.Repaint()
			} else {
				widgets.QMessageBox_Information(nil, "OK", "no Response for Cmd:\n"+serialSendTxt.Text(),
					widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			}
		}
	})

	serialTabPage.SetLayout(serialTabLayout)

	return serialTabPage
}

func getFilesInFolder(root string, ext string) []string {
	var files []string
	log.Printf("looking for files with extension %s in %s\n", root, ext)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		log.Printf("path: %s\n", path)
		if filepath.Ext(path) == ext {
			log.Printf("add %s\n", path)
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		log.Println(file)
	}
	return files
}

func readFileLines(path string) (res []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return res
}

func DecryptData(encarr []byte, key int, size int) []byte {
	arr := make([]byte, size)
	arr = encarr
	for i := 0; i < size; i++ {
		s := int(arr[i])
		t := size + key + i - size/key ^ s
		encarr[i] = byte(t)
	}
	return encarr
}

type nonce struct {
	key    byte
	sector byte
	nt     []byte
	nr     []byte
	ar     []byte
}

func extractNonces(data []byte) (res []nonce) {
	for i := 16; i < (208 - 16); i = i + 16 {
		var n nonce
		n.key = data[i]          //16
		n.sector = data[i+1]     //17
		n.nt = data[i+4 : i+8]   //20-23
		n.nr = data[i+8 : i+12]  //24-27
		n.ar = data[i+12 : i+16] //28-31
		if n.key != byte(0xff) && n.sector != byte(0xff) {
			res = append(res, n)
		}
	}
	return res
}
