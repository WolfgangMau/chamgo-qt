package main

import (
	"github.com/therecipe/qt/widgets"
	"log"
	"os"
	"github.com/WolfgangMau/chamgo-qt/config"
)


var AppName = "Chamgo-QT"
var Cfg config.Config
var Statusbar *widgets.QStatusBar
var DeviceActions config.DeviceActions

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Cfg.Load()
	AppName = Cfg.Gui.Title

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Connected = false

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle(AppName)
	window.SetFixedSize2(950, 600)

	mainlayout := widgets.NewQVBoxLayout()

	tabWidget := widgets.NewQTabWidget(nil)
	tabWidget.AddTab(allSlots(), "Tags")
	tabWidget.AddTab(serialTab(), "Serial")
	tabWidget.SetCurrentIndex(1)

	mainlayout.AddWidget(tabWidget, 0, 0x0020)
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
