package main

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"
	"log"
	"github.com/WolfgangMau/chamgo-qt/eml2dump"
	"encoding/hex"
	"strings"
	"strconv"
	"math"
)

type QTbytes struct {
	LineEdits []*widgets.QLineEdit
	Labels []*widgets.QLabel
}

func dataTab() *widgets.QWidget {
	tablayout:=widgets.NewQHBoxLayout()
	dataTabLayout := widgets.NewQGridLayout(nil)
	dataTabLayout.SetAlign(core.Qt__AlignTop)
	dataTabPage := widgets.NewQWidget(nil, 0)

	dump2emulBtn := widgets.NewQPushButton2("Dump2Emul", nil)
	dump2emulBtn.SetFixedWidth(120)
	dump2emulBtn.ConnectClicked(func(checked bool) {

		fromFileSelect := widgets.NewQFileDialog(nil, 0)
		fromFilename := fromFileSelect.GetOpenFileName(nil, "open File", "", "Bin Files (*.dump *.mfd *.bin);;All Files (*.*)", "", fromFileSelect.Options())
		if fromFilename == "" {
			log.Println("no file selected")
			return
		}

		toFileSelect := widgets.NewQFileDialog(nil, 0)
		toFilename := toFileSelect.GetSaveFileName(nil, "save Data to File", "", "", "", toFileSelect.Options())
		if toFilename == "" {
			log.Println("no file seleted")
			return
		}

		inData,err := eml2dump.File2Bytes(fromFilename)
		if err != nil {
			log.Print(err)
			return
		}

		if !eml2dump.Bytes2Emul(toFilename, inData) {
			log.Print("faiöled to write data to ",toFilename)
			return
		}
	})
	dataTabLayout.AddWidget(dump2emulBtn, 0,0, core.Qt__AlignLeft)


	emul2dumpBtn := widgets.NewQPushButton2("Emul2Dump", nil)
	emul2dumpBtn.SetFixedWidth(120)
	emul2dumpBtn.ConnectClicked(func(checked bool) {

		fromFileSelect := widgets.NewQFileDialog(nil, 0)
		fromFilename := fromFileSelect.GetOpenFileName(nil, "open File", "", "Emulator Files (*.eml *.emul *.txt);;All Files (*.*)", "", fromFileSelect.Options())
		if fromFilename == "" {
			log.Println("no file selected")
			return
		}

		toFileSelect := widgets.NewQFileDialog(nil, 0)
		toFilename := toFileSelect.GetSaveFileName(nil, "save Data to File", "", "", "", toFileSelect.Options())
		if toFilename == "" {
			log.Println("no file seleted")
			return
		}

		filedata, err := eml2dump.File2Bytes(fromFilename)
		if  err != nil {
			log.Print(err)
			return
		}
		bindata,err := hex.DecodeString(strings.Replace(strings.Replace(string(filedata),"\n","",-1),"\r","",-1))
		if err != nil {
			log.Print(err)
			return
		}
		eml2dump.Bytes2File(toFilename,bindata)
	})
	dataTabLayout.AddWidget(emul2dumpBtn, 1,0, core.Qt__AlignLeft)

	tablayout.AddLayout(dataTabLayout,0)
	scrollerA := TagA.Create(true)
	tablayout.AddWidget(scrollerA,1, core.Qt__AlignLeft)
	scrollerB:=TagB.Create(false)
	tablayout.AddWidget(scrollerB,1, core.Qt__AlignLeft)
	//connect(firstScrollbar, SIGNAL(valueChanged(int)), secondScrollbar, SLOT(setValue(int)));
	//connect(secondScrollbar, SIGNAL(valueChanged(int)), firstScrollbar, SLOT(setValue(int)));
	//scrollerA.VerticalScrollBar().ConnectSliderChange(func(scrollerB widgets.QAbstractSlider__SliderChange){
	//})
	dataTabPage.SetLayout(tablayout)
	return dataTabPage
}

func (QTbytesGrid *QTbytes)Create(labelIt bool) *widgets.QScrollArea{
	wrapper := widgets.NewQWidget(nil,0)
	scroller := widgets.NewQScrollArea(nil)
	scroller.SetWidgetResizable(true)
	if labelIt {
		scroller.SetFixedWidth(455)
	} else {
		scroller.SetFixedWidth(400)
	}
	scroller.SetWidget(wrapper)
	sl := widgets.NewQGridLayout(scroller)
	sl.SetSpacing(2)
	sl.SetAlign(core.Qt__AlignLeft)
	startRow := 0
	startcell := 0
	byteCount:=0
	blockCount:=0
	header:=false
	for i:=0;i<=64;i++{
		for i2:=0;i2<=16;i2++ {
			if !header {
				QTbytesGrid.Labels = append(QTbytesGrid.Labels, widgets.NewQLabel(nil,0))
				QTbytesGrid.Labels[i2].SetText(strconv.Itoa(i2))
				sl.AddWidget(QTbytesGrid.Labels[i2], startRow+i, startcell+i2+1, core.Qt__AlignHCenter)

			} else {
				if labelIt && i2==0{
					sector := int(math.Floor(float64(blockCount-1) / 4))
					blockLabel := widgets.NewQLabel2("Block "+strconv.Itoa(blockCount-1), nil, 0)
					blockLabel.SetToolTip("Sector " + strconv.Itoa(sector))
					sl.AddWidget(blockLabel, startRow+i, 0, core.Qt__AlignLeft)
				}
				QTbytesGrid.LineEdits = append(QTbytesGrid.LineEdits, widgets.NewQLineEdit(nil))
				QTbytesGrid.LineEdits[byteCount].SetMaxLength(2)
				QTbytesGrid.LineEdits[byteCount].SetFixedWidth(20)
				QTbytesGrid.LineEdits[byteCount].SetAlignment(core.Qt__AlignHCenter)
				sl.AddWidget(QTbytesGrid.LineEdits[byteCount], startRow+i, startcell+i2+1, core.Qt__AlignLeft)
				byteCount++
			}
		}
		header=true
		blockCount++
	}
	wrapper.SetLayout(sl)
	return scroller
}