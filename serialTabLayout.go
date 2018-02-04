package main

import (
	"github.com/therecipe/qt/widgets"
	"log"
	"strconv"
)

var (
	serialSendButton *widgets.QPushButton
	serialMonitor    *widgets.QPlainTextEdit
	serialPortSelect  *widgets.QComboBox
	deviceSelect *widgets.QComboBox
)

func serialTab() *widgets.QWidget {

	serialTabLayout := widgets.NewQHBoxLayout()
	serialTabPage := widgets.NewQWidget(nil, 0)

	/********************************************** Serial Connect *********************************************/

	serialPorts, _ := getSerialPorts()

	serConLayout := widgets.NewQFormLayout(nil)

	deviceSelect = widgets.NewQComboBox(nil)
	deviceSelect.AddItems(Devices.name)
	deviceSelect.SetCurrentIndex(0)
	deviceSelect.SetFixedWidth(160)
	deviceSelect.SetCurrentIndex(SelectedDeviceId)

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

	serialTabLayout.AddWidget(serialConnectGroup, 0, 0x0020)

	serialConnectButton.ConnectClicked(func(checked bool) {

		Commands.load(deviceSelect.CurrentText())

		if serialConnectButton.Text() == "Connect" {

			err := connectSerial(SerialDevice1)
			if err != nil {
				widgets.QMessageBox_Information(nil, "OK", "can't connect to Serial\n"+string(err.Error()),
					widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				log.Printf("error on connect: %q\n", err)
			} else {
				Device = deviceSelect.CurrentText()
				Commands.load(Device)
				DeviceActions.load(Device)

				//ask for the device-version
				sendSerialCmd(Commands.version + "?")
				if SerialResponse.Code >= 100 {
					serialConnectButton.SetText("Disconnect")
					serialSendButton.SetDisabled(false)
					serialSendButton.Repaint()
					checkForDevices()
				}
				//web got a expected answer from the VERSION(MY) cmd
				if SerialResponse.Code == 101 {
					serialDeviceInfo.SetText("Connected\n" + deviceInfo(SerialResponse.Payload))
					Connected = true
					Statusbar.ShowMessage("Connected to Port: "+serialPortSelect.CurrentText()+" - Device: "+Devices.cdc[SelectedDeviceId]+" - Firmware: "+deviceInfo(SerialResponse.Payload), 0)
					buttonClicked(0)
					buttonClicked(4)
					buttonClicked(1)


				} else {
					widgets.QMessageBox_Information(nil, "OK", "no Version Response from Device!",
						widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				}
			}
		} else {
			err := serialPort.Close()
			if err == nil {
				serialConnectButton.SetText("Connect")
				serialSendButton.SetDisabled(true)
				serialSendButton.Repaint()
				serialDeviceInfo.SetText("not Connected")
				Statusbar.ShowMessage("not Connected", 0)
				Connected = false
				serialPort.Close()
				//GetSlotTicker.Stop()
				//log.Println("GetSlotTicker stopped")
				myProgressBar.zero()
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
