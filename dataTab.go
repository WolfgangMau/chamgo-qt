package main

import (
	"encoding/hex"
	"github.com/WolfgangMau/chamgo-qt/eml2dump"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"log"
	"math"
	"strconv"
	"strings"
)

type QTbytes struct {
	LineEdits []*widgets.QLineEdit
	Labels    []*widgets.QLabel
}

var ScrollLock int

func dataTab() *widgets.QWidget {
	tablayout := widgets.NewQHBoxLayout()


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

		inData, err := eml2dump.File2Bytes(fromFilename)
		if err != nil {
			log.Print(err)
			return
		}

		if !eml2dump.Bytes2Emul(toFilename, inData) {
			log.Print("faiöled to write data to ", toFilename)
			return
		}
	})


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
		if err != nil {
			log.Print(err)
			return
		}
		bindata, err := hex.DecodeString(strings.Replace(strings.Replace(string(filedata), "\n", "", -1), "\r", "", -1))
		if err != nil {
			log.Print(err)
			return
		}
		eml2dump.Bytes2File(toFilename, bindata)
	})


	loadTagABtn := widgets.NewQPushButton2("Load TagA", nil)
	loadTagABtn.SetFixedWidth(120)
	loadTagABtn.ConnectClicked(func(checked bool) {
		fromFileSelect := widgets.NewQFileDialog(nil, 0)
		fromFilename := fromFileSelect.GetOpenFileName(nil, "open File", "", "Bin Files (*.dump *.mfd *.bin);;All Files (*.*)", "", fromFileSelect.Options())
		if fromFilename == "" {
			log.Println("no file selected")
			return
		}
		TagA.FillFromFile(fromFilename)
	})

	loadTagBBtn := widgets.NewQPushButton2("Load TagB", nil)
	loadTagBBtn.SetFixedWidth(120)
	loadTagBBtn.ConnectClicked(func(checked bool) {
		fromFileSelect := widgets.NewQFileDialog(nil, 0)
		fromFilename := fromFileSelect.GetOpenFileName(nil, "open File", "", "Bin Files (*.dump *.mfd *.bin);;All Files (*.*)", "", fromFileSelect.Options())
		if fromFilename == "" {
			log.Println("no file selected")
			return
		}
		TagB.FillFromFile(fromFilename)
	})


	diffTagBtn := widgets.NewQPushButton2("Diff A/B", nil)
	diffTagBtn.SetFixedWidth(120)
	diffTagBtn.ConnectClicked(func(checked bool) {
		DiffTags(TagA, TagB)
	})


	mapTagBtn := widgets.NewQPushButton2("Map Mf1K", nil)
	mapTagBtn.SetFixedWidth(120)
	mapTagBtn.ConnectClicked(func(checked bool) {
		TagA.Map("Mf1K")
		TagB.Map("Mf1K")
	})

	var slotItems  []string
	slotSelect := widgets.NewQComboBox(nil)
	slotSelect.SetFixedWidth(120)
	slotSelect.ConnectCurrentIndexChanged(func(index int) {
		if Connected && index>0 {
			data := getSlotBytes(index - 1)
			log.Print("data len: ", len(data))
			if len(data)>0 {
				TagA.FillFromBytes(data)
			}
		}
	})
	slotItems = append(slotItems, "Load from Slot")
	for i:=Cfg.Device[SelectedDeviceId].Config.Slot.First; i<=Cfg.Device[SelectedDeviceId].Config.Slot.Last; i++ {
		slotItems = append(slotItems, "Slot "+strconv.Itoa(i))
	}
	slotSelect.AddItems(slotItems)


	lockScrollChk := widgets.NewQCheckBox2("ScrollLock", nil)
	lockScrollChk.ConnectStateChanged(func(state int) {

			ScrollLock = state
			log.Printf("ScrollLock: %d", ScrollLock)

	})

	//left menuy
	dataTabLayout := widgets.NewQGridLayout(nil)
	dataTabLayout.SetAlign(core.Qt__AlignTop)

	fileGroup := widgets.NewQGroupBox2("File based",nil)
	fileLayout := widgets.NewQVBoxLayout()

	fileLayout.AddWidget(dump2emulBtn, 0,core.Qt__AlignLeft)
	fileLayout.AddWidget(emul2dumpBtn, 0,core.Qt__AlignLeft)
	fileLayout.AddWidget(loadTagABtn, 0,core.Qt__AlignLeft)
	fileLayout.AddWidget(loadTagBBtn, 0,core.Qt__AlignLeft)

	fileGroup.SetLayout(fileLayout)
	dataTabLayout.AddWidget(fileGroup, 0, 0, core.Qt__AlignLeft)

	//dataTabLayout.AddWidget(dump2emulBtn, 0, 0, core.Qt__AlignLeft)
	//dataTabLayout.AddWidget(emul2dumpBtn, 1, 0, core.Qt__AlignLeft)
	//dataTabLayout.AddWidget(loadTagABtn, 2, 0, core.Qt__AlignLeft)
	//dataTabLayout.AddWidget(loadTagBBtn, 3, 0, core.Qt__AlignLeft)

	dataTabLayout.AddWidget(diffTagBtn, 1, 0, core.Qt__AlignLeft)
	dataTabLayout.AddWidget(mapTagBtn, 2, 0, core.Qt__AlignLeft)
	dataTabLayout.AddWidget(lockScrollChk, 3, 0, core.Qt__AlignLeft)
	dataTabLayout.AddWidget(slotSelect, 4, 0, core.Qt__AlignLeft)
	tablayout.AddLayout(dataTabLayout, 0)

	scrollerA := TagA.Create(true)
	tablayout.AddWidget(scrollerA, 1, core.Qt__AlignLeft)

	scrollerB := TagB.Create(false)
	tablayout.AddWidget(scrollerB, 1, core.Qt__AlignLeft)

	scrollerA.VerticalScrollBar().ConnectValueChanged(func(positionA int) {
		if ScrollLock == 2 {
			scrollerB.VerticalScrollBar().SetValue(positionA)
		}
	})
	scrollerB.VerticalScrollBar().ConnectValueChanged(func(positionB int) {
		if ScrollLock == 2 {
			scrollerA.VerticalScrollBar().SetValue(positionB)
		}
	})
	dataTabPage.SetLayout(tablayout)
	return dataTabPage
}

func (QTbytesGrid *QTbytes) Create(labelIt bool) *widgets.QScrollArea {
	wrapper := widgets.NewQWidget(nil, 0)
	scroller := widgets.NewQScrollArea(nil)
	scroller.SetWidgetResizable(true)
	if labelIt {
		scroller.SetFixedWidth(435)
	} else {
		scroller.SetFixedWidth(380)
	}
	scroller.SetWidget(wrapper)
	sl := widgets.NewQGridLayout(scroller)
	sl.SetSpacing(2)
	sl.SetAlign(core.Qt__AlignLeft)
	startRow := 0
	startcell := 0
	byteCount := 0
	blockCount := 0
	header := false
	for i := 0; i <= 64; i++ {
		for i2 := 0; i2 <= 15; i2++ {
			if !header {
				QTbytesGrid.Labels = append(QTbytesGrid.Labels, widgets.NewQLabel(nil, 0))
				QTbytesGrid.Labels[i2].SetText(strconv.Itoa(i2))
				sl.AddWidget(QTbytesGrid.Labels[i2], startRow+i, startcell+i2+1, core.Qt__AlignHCenter)

			} else {
				if labelIt && i2 == 0 {
					sector := int(math.Floor(float64(blockCount-1) / 4))
					blockLabel := widgets.NewQLabel2("Block "+strconv.Itoa(blockCount-1), nil, 0)
					blockLabel.SetToolTip("Sector " + strconv.Itoa(sector))
					sl.AddWidget(blockLabel, startRow+i, 0, core.Qt__AlignLeft)
				}
				QTbytesGrid.LineEdits = append(QTbytesGrid.LineEdits, widgets.NewQLineEdit(nil))
				QTbytesGrid.LineEdits[byteCount].SetToolTip("Byte #"+strconv.Itoa(byteCount))
				QTbytesGrid.LineEdits[byteCount].SetMaxLength(2)
				QTbytesGrid.LineEdits[byteCount].SetFixedWidth(20)
				QTbytesGrid.LineEdits[byteCount].SetAlignment(core.Qt__AlignHCenter)
				sl.AddWidget(QTbytesGrid.LineEdits[byteCount], startRow+i, startcell+i2+1, core.Qt__AlignLeft)
				byteCount++
			}
		}
		header = true
		blockCount++
	}
	wrapper.SetLayout(sl)
	return scroller
}

func (QTbytesGrid *QTbytes) FillFromFile(filename string) {
	data, _ := eml2dump.File2Bytes(filename)
	if len(QTbytesGrid.LineEdits) != len(data) {
		log.Printf("data-Len missmamatch grid: %d - file: %d", len(QTbytesGrid.LineEdits), len(data))
		return
	}
	for i, b := range data {
		QTbytesGrid.LineEdits[i].SetText(hex.EncodeToString([]byte{b}))
	}
}

func (QTbytesGrid *QTbytes) FillFromBytes(data []byte) {
	if len(QTbytesGrid.LineEdits) != len(data) {
		log.Printf("data-Len missmamatch grid: %d - data: %d", len(QTbytesGrid.LineEdits), len(data))
		return
	}
	for i, b := range data {
		QTbytesGrid.LineEdits[i].SetText(hex.EncodeToString([]byte{b}))
	}
}

func (QTbytesGrid *QTbytes) SetColor(index int, color string, selcolor string) {
	style := "background: " + rgbColorString(color) + " selection-background-color: " + rgbColorString(selcolor)
	QTbytesGrid.LineEdits[index].SetStyleSheet(style)
}

func rgbColorString(color string) (res string) {
	switch color {
	case "red":
		res = "rgb(255, 0, 0);"
	case "lightred":
		res = "rgb(255, 71, 26);"

	case "green":
		res = "rgb(51, 153, 51);"
	case "lightgreen":
		res = "rgb(0, 204, 0);"
	case "limegreen":
		res = "rgb(0, 255, 0);"

	case "blue":
		res = "rgb(0, 102, 255);"
	case "lightblue":
		res = "rgb(51, 204, 255);"

	case "yellow":
		res = "rgb(255, 255, 0);"

	case "purple":
		res = "rgb(153, 51, 255);"

	case "magenta":
		res = "rgb(255, 0, 255);"

	case "grey":
		res = "rgb(179, 179, 179);"

	case "defaultsel":
		res = "rgb(51, 204, 255);"

	default:
		res = "rgb(255,255,255);"
	}
	return res
}

func DiffTags(tagA QTbytes, tagB QTbytes) {
	for i, le := range tagA.LineEdits {
		if strings.ToLower(le.Text()) != strings.ToLower(tagB.LineEdits[i].Text()) {
			tagA.SetColor(i, "lightgreen", "defaultsel")
			tagB.SetColor(i, "lightred", "defaultsel")
		} else {
			tagA.SetColor(i, "default", "defaultsel")
			tagB.SetColor(i, "default", "defaultsel")
		}
	}
}

func (tag *QTbytes)Map(mapName string) {
	sectorblock:=0
	blockbyte:=0
	sectror:=0

	for i,_ := range tag.LineEdits {
		if i>15 && float64(i%16)==0.0 {
			sectorblock++
			if float64((sectorblock)%4)==0.0 {
				sectror++
				sectorblock=0
			}
			blockbyte=0
		}
		tooltip:="Sector #"+strconv.Itoa(sectror)+" - sectorBlock #"+strconv.Itoa(sectorblock)+" - blockByte #"+strconv.Itoa(blockbyte)
		if mapName == "Mf1K" {
			switch {
			// UID
			case i <= 3:
				tag.SetColor(i, "magenta", "defaultsel")
				tooltip += "\nUID"
				// BCC
			case i == 4:
				tag.SetColor(i, "yellow", "defaultsel")
				tooltip += "\nBCC"
				// SAK
			case i == 5:
				tag.SetColor(i, "magenta", "defaultsel")
				tooltip += "\nSAK"
				// SAK
			case i == 6:
				tag.SetColor(i, "magenta", "defaultsel")
				tooltip += "\nATQA_0"
				// SAK
			case i == 7:
				tag.SetColor(i, "magenta", "defaultsel")
				// SAK
			case i >7 && i< 16:
				tag.SetColor(i, "grey", "defaultsel")
				tooltip += "\nManufacuter Data"
				// KEYA/B
			case sectorblock == 3 && (blockbyte <= 5 || blockbyte > 9):
				tag.SetColor(i, "magenta", "defaultsel")
				if blockbyte <= 5 {
					tooltip += "\nKEY A"
				} else {
					tooltip += "\nKEY B"
				}
				// Permissions
			case sectorblock == 3 && blockbyte > 5 && blockbyte < 9:
				tag.SetColor(i, "grey", "defaultsel")
				tooltip += "\nPERM"
				// Peneral Purpose
			case sectorblock == 3 && blockbyte == 9:
				tag.SetColor(i, "lightblue", "yellow")
				tooltip += "\nGPS"
			}
		}
		tag.LineEdits[i].SetToolTip(tooltip)
		blockbyte++
	}
}