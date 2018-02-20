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
	if _, err := getSerialPorts(); err != nil {
		log.Println(err)
	}
	Cfg.Load()
	dn := Cfg.Device[SelectedDeviceId].Name
	DeviceActions.Load(Cfg.Device[SelectedDeviceId].CmdSet, dn)
}

func main() {
	var f *os.File
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.OpenFile(config.Apppath()  + string(os.PathSeparator) + "chamgo-qt.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
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
