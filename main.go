package main

import (
	"github.com/therecipe/qt/widgets"
	"log"
	"os"
)

var AppName = "Chamgo-QT"
var Connected bool
var Statusbar *widgets.QStatusBar

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Connected = false
	ActionButtons = []string{"Select All", "Select None", "Apply", "Clear", "Refresh", "Set Active", "mfkey32", "Upload", "Download"}

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle(AppName)
	window.SetFixedSize2(950, 600)

	mainlayout := widgets.NewQVBoxLayout()

	tabWidget := widgets.NewQTabWidget(nil)
	tabWidget.AddTab(allSlots(), "Tags")
	tabWidget.AddTab(serialTab(), "Serial")
	//tabWidget.SetFixedSize2(800,500)
	tabWidget.SetCurrentIndex(1)

	mainlayout.AddWidget(tabWidget, 0, 0x0020)
	mainlayout.SetAlign(33)

	myProgressBar.widget = widgets.NewQProgressBar(window)
	myProgressBar.widget.SetRange(0, 100)
	myProgressBar.widget.SetVisible(true)
	myProgressBar.widget.ShowDefault()
	mainlayout.AddWidget(myProgressBar.widget, 0, 0x0020)

	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetLayout(mainlayout)
	window.SetCentralWidget(mainWidget)

	Statusbar = widgets.NewQStatusBar(window)
	Statusbar.ShowMessage("not Connected", 0)
	window.SetStatusBar(Statusbar)

	// Show the window
	window.Show()

	// Execute app
	app.Exec()
}
