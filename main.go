package main

import (
	"github.com/WolfgangMau/chamgo-qt/config"
	"github.com/therecipe/qt/widgets"
	"log"
	"os"
)

//Global Variables - StateStorage
var AppName = "Chamgo-QT"
var Cfg config.Config
var Statusbar *widgets.QStatusBar
var DeviceActions config.DeviceActions
var MyTabs *widgets.QTabWidget


func initcfg() {
	if _,err := getSerialPorts(); err != nil {
		log.Println(err)
	}
	Cfg.Load()
	dn := Cfg.Device[SelectedDeviceId].Name
	DeviceActions.Load(Cfg.Device[SelectedDeviceId].CmdSet, dn)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	initcfg()
	AppName = Cfg.Gui.Title

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Connected = false

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle(AppName)
	window.SetFixedSize2(950, 600)

	mainlayout := widgets.NewQVBoxLayout()

	MyTabs = widgets.NewQTabWidget(nil)
	MyTabs.AddTab(allSlots(), "Tags")
	MyTabs.AddTab(serialTab(), "Serial")
	MyTabs.SetCurrentIndex(1)

	mainlayout.AddWidget(MyTabs, 0, 0x0020)
	mainlayout.SetAlign(33)

	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetLayout(mainlayout)
	window.SetCentralWidget(mainWidget)

	Statusbar = widgets.NewQStatusBar(window)
	Statusbar.ShowMessage("not Connected", 0)
	window.SetStatusBar(Statusbar)

	checkForDevices()
	// Show the window
	window.Show()

	// Execute app
	app.Exec()
}
