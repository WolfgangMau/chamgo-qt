package main

import (
	"encoding/hex"
	"github.com/WolfgangMau/chamgo-qt/config"
	"github.com/WolfgangMau/chamgo-qt/crc16"
	"github.com/WolfgangMau/chamgo-qt/eml2dump"
	"github.com/WolfgangMau/chamgo-qt/xmodem"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type QTbytes struct {
	LineEdits []*widgets.QLineEdit
	Labels    []*widgets.QLabel
}

var ScrollLock int
var RadioTag string

func dataTab() *widgets.QWidget {
	tablayout := widgets.NewQHBoxLayout()
	dataTabPage := widgets.NewQWidget(nil, 0)
	dump2emulBtn := widgets.NewQPushButton2("Dump > Emul", nil)
	dump2emulBtn.SetToolTip("convert Binary-File to ASCII-File")
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

	emul2dumpBtn := widgets.NewQPushButton2("Emul > Dump", nil)
	emul2dumpBtn.SetToolTip("convert ASCII-File to Binary-File")
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

	loadTagABtn := widgets.NewQPushButton2("Dump > Tag A", nil)
	loadTagABtn.SetToolTip("load Data from Binary-File to Tag A")
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

	loadTagBBtn := widgets.NewQPushButton2("Dump > Tag B", nil)
	loadTagBBtn.SetToolTip("load Data from Binary-File to Tag A+B")
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

	diffTagBtn := widgets.NewQPushButton2("Diff Tags", nil)
	diffTagBtn.SetToolTip("show differences between Tag A & Tag B")
	diffTagBtn.SetFixedWidth(120)
	diffTagBtn.ConnectClicked(func(checked bool) {
		DiffTags(TagA, TagB)
	})

	mapfiles := config.GetFilesInFolder(config.Apppath()+string(filepath.Separator)+runtime.GOOS+string(filepath.Separator)+"maps"+string(filepath.Separator), ".map")
	mapSelect := widgets.NewQComboBox(nil)
	mapSelect.SetFixedWidth(120)
	mapSelect.AddItems([]string{"mappings"})
	if len(mapfiles) > 0 {
		mapSelect.AddItems(mapfiles)
	} else {
		log.Print("creating dummy-map dummy.map")
		temp := config.DefaultMap
		temp.Save("dummy.map")
		mapSelect.AddItems([]string{"dummy.map"})
	}
	mapSelect.ConnectCurrentIndexChanged(func(index int) {
		if index > 0 {
			tagmap := config.TagMap{}
			tagmap.Load(mapSelect.CurrentText())
			TagA.Map(tagmap)
			TagB.Map(tagmap)
			mapSelect.SetCurrentIndex(0)
			mapSelect.Repaint()
		}
	})

	var slotItems []string
	slotSelect := widgets.NewQComboBox(nil)
	slotSelect.SetToolTip("load Data from Slot to Tag A or B")
	slotSelect.SetFixedWidth(120)
	slotSelect.ConnectCurrentIndexChanged(func(index int) {
		if Connected && index > 0 {
			data := getSlotBytes(index - 1)
			if len(data) > 0 {
				if RadioTag == "A" {
					TagA.FillFromBytes(data)
				} else {
					TagB.FillFromBytes(data)
				}
				slotSelect.SetCurrentIndex(0)
			}
		}
	})
	slotSelect.AddItems([]string{"Slot > Tag"})

	for i := Cfg.Device[SelectedDeviceId].Config.Slot.First; i <= Cfg.Device[SelectedDeviceId].Config.Slot.Last; i++ {
		slotItems = append(slotItems, "Slot "+strconv.Itoa(i))
	}
	slotSelect.AddItems(slotItems)

	slot2Select := widgets.NewQComboBox(nil)
	slot2Select.SetToolTip("load Data from Slot to Tag A or B")
	slot2Select.SetFixedWidth(120)
	slot2Select.ConnectCurrentIndexChanged(func(index int) {
		if Connected && index > 0 {
			var data []byte
			var tag QTbytes
			switch RadioTag {
			case "A":
				tag = TagA
			case "B":
				tag = TagB
			}
			for _, le := range tag.LineEdits {
				lbyte, err := hex.DecodeString(le.Text())

				if err != nil {
					log.Print(err)
				} else {
					if len(lbyte) == 1 {
						data = append(data, lbyte[0])
					}
				}
			}
			if len(data) > 0 {
				if RadioTag == "A" {
					TagA.Tag2Slot(index-1, data)
				} else {
					TagB.Tag2Slot(index-1, data)
				}
				slot2Select.SetCurrentIndex(0)
			}
		}
	})
	slot2Select.AddItems([]string{"Tag > Slot"})
	slot2Select.AddItems(slotItems)

	lockScrollChk := widgets.NewQCheckBox2("ScrollLock", nil)
	lockScrollChk.SetToolTip("Scroll together or separately")
	lockScrollChk.SetChecked(true)
	ScrollLock = 2
	lockScrollChk.ConnectStateChanged(func(state int) {

		ScrollLock = state
		log.Printf("ScrollLock: %d", ScrollLock)

	})

	//left menuy

	dataTabLayout := widgets.NewQGridLayout(nil)
	dataTabLayout.SetAlign(core.Qt__AlignTop)

	//slot action
	slotGroup := widgets.NewQGroupBox2("Slot Import / Export", nil)
	slotlayout := widgets.NewQVBoxLayout()
	radioLayout := widgets.NewQHBoxLayout()
	tagARadio := widgets.NewQRadioButton2("Tag A", slotGroup)
	tagARadio.ConnectClicked(func(checked bool) {
		if checked {
			RadioTag = "A"
		} else {
			RadioTag = "B"
		}
	})
	tagARadio.SetChecked(true)
	RadioTag = "A"
	tagBRadio := widgets.NewQRadioButton2("Tag B", slotGroup)
	tagBRadio.ConnectClicked(func(checked bool) {
		if checked {
			RadioTag = "B"
		} else {
			RadioTag = "A"
		}
	})
	radioLayout.AddWidget(tagARadio, 0, core.Qt__AlignCenter)
	radioLayout.AddWidget(tagBRadio, 0, core.Qt__AlignCenter)
	slotlayout.AddLayout(radioLayout, 0)
	slotlayout.AddWidget(slotSelect, 0, core.Qt__AlignCenter)
	slotlayout.AddWidget(slot2Select, 0, core.Qt__AlignCenter)
	slotGroup.SetFixedWidth(155)
	slotGroup.SetLayout(slotlayout)

	//File based action
	fileGroup := widgets.NewQGroupBox2("File Import / Export", nil)
	fileLayout := widgets.NewQVBoxLayout()
	fileLayout.AddWidget(dump2emulBtn, 0, core.Qt__AlignCenter)
	fileLayout.AddWidget(emul2dumpBtn, 0, core.Qt__AlignCenter)
	fileLayout.AddWidget(loadTagABtn, 0, core.Qt__AlignCenter)
	fileLayout.AddWidget(loadTagBBtn, 0, core.Qt__AlignCenter)
	fileGroup.SetFixedWidth(155)
	fileGroup.SetLayout(fileLayout)

	//Diff/Map based actions
	diffGroup := widgets.NewQGroupBox2("Diff / Map", nil)
	difflayout := widgets.NewQVBoxLayout()
	difflayout.AddWidget(diffTagBtn, 0, core.Qt__AlignCenter)
	difflayout.AddWidget(mapSelect, 0, core.Qt__AlignCenter)
	difflayout.AddWidget(lockScrollChk, 0, core.Qt__AlignCenter)
	diffGroup.SetFixedWidth(155)
	diffGroup.SetLayout(difflayout)

	dataTabLayout.AddWidget(fileGroup, 0, 0, core.Qt__AlignCenter)
	dataTabLayout.AddWidget(slotGroup, 1, 0, core.Qt__AlignCenter)
	dataTabLayout.AddWidget(diffGroup, 2, 0, core.Qt__AlignCenter)
	tablayout.AddLayout(dataTabLayout, 0)

	tagALayout := widgets.NewQVBoxLayout()
	scrollerA := TagA.Create(true)
	tagAInfo := widgets.NewQLabel2("Tag A", nil, 0)
	tagALayout.AddWidget(tagAInfo, 0, core.Qt__AlignCenter)
	tagALayout.AddWidget(scrollerA, 0, core.Qt__AlignLeft)
	tablayout.AddLayout(tagALayout, 1)

	tagBLayout := widgets.NewQVBoxLayout()
	scrollerB := TagB.Create(false)
	tagBInfo := widgets.NewQLabel2("Tag B", nil, 0)
	tagBLayout.AddWidget(tagBInfo, 0, core.Qt__AlignCenter)
	tagBLayout.AddWidget(scrollerB, 0, core.Qt__AlignLeft)
	tablayout.AddLayout(tagBLayout, 0)

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
				QTbytesGrid.LineEdits[byteCount].SetToolTip("Byte #" + strconv.Itoa(byteCount))
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

func (QTbytesGrid *QTbytes) Tag2Slot(slotnumber int, data []byte) {
	hardwareSlot := slotnumber + Cfg.Device[SelectedDeviceId].Config.Slot.Offset
	sendSerialCmd(DeviceActions.SelectSlot + strconv.Itoa(hardwareSlot))
	p := Bytes2Packets(data)
	//set chameleon into receiver-mode
	sendSerialCmd(DeviceActions.StartUpload)
	if SerialResponse.Code == 110 {
		//start uploading packets
		xmodem.Send(SerialPort, p)
	}
}

func (QTbytesGrid *QTbytes) SetColor(index int, color string, selcolor string) {
	style := "background: " + rgbColorString(color) + " selection-background-color: " + rgbColorString(selcolor)
	QTbytesGrid.LineEdits[index].SetStyleSheet(style)
	QTbytesGrid.LineEdits[index].Repaint()
}

func (QTbytesGrid *QTbytes) SetTooltip(index int, tip string) {
	QTbytesGrid.LineEdits[index].SetToolTip(tip)
	QTbytesGrid.LineEdits[index].Repaint()
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

func (tag *QTbytes) Map(tm config.TagMap) {
	if len(tm.Mappings) > 0 {
		for _, m := range tm.Mappings {
			if len(m.MapBytes) > 0 {
				if len(m.MapFuncs) == 0 {
					for _, mb := range m.MapBytes {
						tag.SetColor(mb.Pos, mb.Color, "defaultsel")
						tag.SetTooltip(mb.Pos, mb.Tooltip)
					}
				} else {
					words := strings.Split(m.MapFuncs[0].Name, " ")
					if len(words) == 3 {
						switch strings.ToLower(words[0]) {
						case "each":
							switch strings.ToLower(words[1]) {
							case "sectorblock":
								targetBlock, _ := strconv.Atoi(words[2])
								sectorBlock := 0
								sector := 0
								blockByte := 0
								for i := range tag.LineEdits {
									if blockByte == 16 {
										sectorBlock++
										blockByte = 0
										if sectorBlock == 4 {
											sector++
											sectorBlock = 0
										}
									}
									if sectorBlock == targetBlock {
										if blockByte >= m.Start && blockByte <= m.End {
											tag.SetColor(i, m.MapBytes[0].Color, "defaultsel")
											tag.SetTooltip(i, "Sector: "+strconv.Itoa(sector)+" - sectorBlock: "+strconv.Itoa(sectorBlock)+" blockByte: "+strconv.Itoa(i)+"\n"+m.MapBytes[0].Tooltip)
										}
									}
									blockByte++
								}
							}
						case "check":
							switch strings.ToLower(words[1]) {
							case "bcc":
								if tag.CalcBCC(4) == tag.LineEdits[4].Text() {
									tag.SetColor(4, words[2], "defaultsel")
									tag.SetTooltip(4, "BCC OK")
								} else {
									tag.SetColor(4, m.MapBytes[0].Color, "defaultsel")
									tag.SetTooltip(4, "BCC != "+tag.CalcBCC(4))
								}
							}
						}
					}
				}
			}
		}
	}
}

func (tag *QTbytes) GetId(size int) string {
	var res string
	for i := 0; i < size; i++ {
		res += tag.LineEdits[i].Text()
	}
	return res
}

func (tag *QTbytes) CalcBCC(size int) string {
	id := tag.GetId(size)
	b, _ := hex.DecodeString(id)
	bcc2 := hex.EncodeToString([]byte{crc16.GetBCC(b)})
	return bcc2
}
