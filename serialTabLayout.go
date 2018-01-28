package main

import (
	"github.com/therecipe/qt/widgets"
	"log"
)

var (
	serialSendButton   *widgets.QPushButton
	serialMonitor      *widgets.QPlainTextEdit
	serialResponseList []string
	serialResponseStr  string
)

func serialTab() *widgets.QWidget {

	Devices.load()

	serialTabLayout := widgets.NewQHBoxLayout()
	serialTabPage := widgets.NewQWidget(nil, 0)

	/********************************************** Serial Connect *********************************************/

	serialPorts, _ := getSerialPorts()

	serConLayout := widgets.NewQFormLayout(nil)

	deviceSelect := widgets.NewQComboBox(nil)
	deviceSelect.AddItems(Devices.name)
	deviceSelect.SetCurrentIndex(0)
	deviceSelect.SetFixedWidth(160)

	serialPortSelect := widgets.NewQComboBox(nil)
	serialPortSelect.AddItems(serialPorts)
	serialPortSelect.SetFixedWidth(160)

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
			err := connectSerial(serialPortSelect.CurrentText())
			if err != nil {
				log.Fatal("can't connect to Serial: ", err)
			} else {
				Device = deviceSelect.CurrentText()
				Commands.load(Device)
				DeviceActions.load(Device)

				serialConnectButton.SetText("Disconnect")
				serialSendButton.SetDisabled(false)
				serialSendButton.Repaint()
				serialResponseList = serialCMD(Commands.version + "?")

				if len(serialResponseList) == 1 {
					serialDeviceInfo.SetText("Connected\n" + deviceInfo(serialResponseList[0]))
					Connected = true
					Statusbar.ShowMessage("Connected to Port: "+serialPortSelect.CurrentText()+" - Device: "+Device+" - Firmware: "+deviceInfo(serialResponseList[0]), 0)
					populateSlots()
					checkCurrentSelection()

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
				GetSlotTicker.Stop()
				log.Println("GetSlotTicker stopped")
			}
		}

	})

	/********************************************** Serial Monitor *********************************************/

	serMonitorLayout := widgets.NewQVBoxLayout()
	serSendLayout := widgets.NewQHBoxLayout()

	serialMonitor = widgets.NewQPlainTextEdit(nil)
	serialMonitor.AppendPlainText("")
	serialMonitor.SetFixedHeight(400)
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
			serialResponseStr = sendSerial(serialSendTxt.Text())
			if serialResponseStr != "" {
				serialMonitor.AppendPlainText("-> " + serialSendTxt.Text())
				serialMonitor.AppendPlainText("<- " + serialResponseStr)
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
