package main

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"log"
	"strconv"
)

type Slot struct {
	widgets.QMainWindow
	slotl *widgets.QLabel
	slot  *widgets.QCheckBox
	model *widgets.QLabel
	mode  *widgets.QComboBox
	uidl  *widgets.QLabel
	uid   *widgets.QLineEdit
	btnsl *widgets.QLabel
	btns  *widgets.QComboBox
	btnll *widgets.QLabel
	btnl  *widgets.QComboBox
	sizel *widgets.QLabel
	size  *widgets.QLineEdit
}

type SlotHLayout struct {
	l *widgets.QHBoxLayout
}
type SlotVLayout struct {
	l *widgets.QVBoxLayout
}

type SlotBox struct {
	b *widgets.QGroupBox
}

type AButton struct {
	b *widgets.QPushButton
}

var AButtons [9]AButton

var Slots [8]Slot
var Slotlayouts [2]SlotHLayout
var SlotHlayouts [8]SlotHLayout
var SlotGroupVlayouts [8]SlotVLayout
var Slotboxes [8]SlotBox
var ActionButtons []string

func allSlots() *widgets.QWidget {
	bold := gui.NewQFont()
	bold.SetBold(true)

	slotsTabLayout := widgets.NewQGridLayout(nil)
	slotsTabPage := widgets.NewQWidget(nil, 0)
	var c = 0
	for i := 0; i <= 1; i++ {
		Slotlayouts[i].l = widgets.NewQHBoxLayout()
		Slotlayouts[i].l.SetAlign(33)

		//SlotGroupVlayouts[i].l = widgets.NewQVBoxLayout()
		//SlotGroupVlayouts[i].l.SetAlign(33)

		//var gc=0
		for s := 0; s <= 3; s++ {
			/************* Slot checkbox ************/
			SlotHlayouts[i].l = widgets.NewQHBoxLayout()
			SlotHlayouts[i].l.SetContentsMargins(10, 0, 0, 0)
			SlotHlayouts[i].l.Stretch(1)

			log.Printf("building Slot %d", c)
			Slots[c].slotl = widgets.NewQLabel(nil, 0)

			Slots[c].slotl.SetText("Slot " + strconv.Itoa(c+1))
			Slots[c].slotl.SetContentsMargins(10, 0, 0, 0)
			Slots[c].slotl.SetFont(bold)

			Slots[c].slot = widgets.NewQCheckBox(nil)
			Slots[c].slot.SetChecked(false)

			SlotHlayouts[i].l.AddWidget(Slots[c].slot, 0, 0x0001)
			SlotHlayouts[i].l.AddWidget(Slots[c].slotl, 1, 0x0001)

			/************* Slot Group ************/
			boxlayout := widgets.NewQGridLayout(nil)
			boxlayout.SetAlign(33)

			Slots[c].model = widgets.NewQLabel(nil, 0)
			Slots[c].model.SetText("Mode")
			Slots[c].mode = widgets.NewQComboBox(nil)
			Slots[c].mode.SetFixedWidth(127)
			boxlayout.AddWidget(Slots[c].model, 0, 0, 0x0001)
			boxlayout.AddWidget(Slots[c].mode, 0, 1, 0x0001)

			Slots[c].uidl = widgets.NewQLabel(nil, 0)
			Slots[c].uidl.SetText("UID")
			Slots[c].uid = widgets.NewQLineEdit(nil)
			Slots[c].uid.SetFixedWidth(121)
			boxlayout.AddWidget(Slots[c].uidl, 1, 0, 0x0001)
			boxlayout.AddWidget(Slots[c].uid, 1, 1, 0x0001)

			Slots[c].btnsl = widgets.NewQLabel(nil, 0)
			Slots[c].btnsl.SetText("Btn Short")
			Slots[c].btns = widgets.NewQComboBox(nil)
			Slots[c].btns.SetFixedWidth(127)
			boxlayout.AddWidget(Slots[c].btnsl, 2, 0, 0x0001)
			boxlayout.AddWidget(Slots[c].btns, 2, 1, 0x0001)

			Slots[c].btnll = widgets.NewQLabel(nil, 0)
			Slots[c].btnll.SetText("Btn Long")
			Slots[c].btnl = widgets.NewQComboBox(nil)
			Slots[c].btnl.SetFixedWidth(127)
			boxlayout.AddWidget(Slots[c].btnll, 3, 0, 0x0001)
			boxlayout.AddWidget(Slots[c].btnl, 3, 1, 0x0001)

			Slots[c].sizel = widgets.NewQLabel(nil, 0)
			Slots[c].sizel.SetText("Size")
			Slots[c].size = widgets.NewQLineEdit(nil)
			Slots[c].size.SetDisabled(true)
			Slots[c].size.SetFixedWidth(121)
			boxlayout.AddWidget(Slots[c].sizel, 4, 0, 0x0001)
			boxlayout.AddWidget(Slots[c].size, 4, 1, 0x0001)

			SlotGrouplayout := widgets.NewQVBoxLayout()
			SlotGrouplayout.AddLayout(boxlayout, 0)

			Slotboxes[i].b = widgets.NewQGroupBox(nil)
			Slotboxes[i].b.SetLayout(SlotGrouplayout)

			SlotGroupVlayouts[i].l = widgets.NewQVBoxLayout()
			SlotGroupVlayouts[i].l.SetSpacing(0)
			SlotGroupVlayouts[i].l.AddLayout(SlotHlayouts[i].l, 0)
			SlotGroupVlayouts[i].l.AddWidget(Slotboxes[i].b, 1, 0x0001)

			slotsTabLayout.AddLayout(SlotGroupVlayouts[i].l, i, s, 0x0020)

			c++
		}
	}

	ActionButtons = []string{"Select All", "Select None", "Apply", "Clear", "Refresh", "Set Active", "mfkey32", "Upload", "Download"}
	abtnLayout := widgets.NewQGridLayout(nil)
	for i, s := range ActionButtons {
		AButtons[i].b = widgets.NewQPushButton2(s, nil)
		abtnLayout.AddWidget(AButtons[i].b, 0, i, 0x0004)
	}
	AButtonGroup := widgets.NewQGroupBox2("Available Actions", nil)
	AButtonGroup.SetLayout(abtnLayout)
	A2ButtonLayout := widgets.NewQHBoxLayout()
	A2ButtonLayout.AddWidget(AButtonGroup, 1, 0x0001)
	slotsTabLayout.AddLayout2(A2ButtonLayout, 3, 0, 1, 5, 0x0001)

	//for i:=0; i<len(Slots); i++ {
	//	Slots[c].slot.ConnectStateChanged(func(checked int) {
	//		log.Printf("Checked %d - state: %d\n", i, checked)
	//	})
	//}

	//
	Slots[0].slot.ConnectStateChanged(func(checked int) {
		slotChecked(0, checked)
	})
	Slots[1].slot.ConnectStateChanged(func(checked int) {
		slotChecked(1, checked)
	})
	Slots[2].slot.ConnectStateChanged(func(checked int) {
		slotChecked(2, checked)
	})
	Slots[3].slot.ConnectStateChanged(func(checked int) {
		slotChecked(3, checked)
	})
	Slots[4].slot.ConnectStateChanged(func(checked int) {
		slotChecked(4, checked)
	})
	Slots[5].slot.ConnectStateChanged(func(checked int) {
		slotChecked(5, checked)
	})
	Slots[6].slot.ConnectStateChanged(func(checked int) {
		slotChecked(6, checked)
	})
	Slots[7].slot.ConnectStateChanged(func(checked int) {
		slotChecked(7, checked)
	})

	AButtons[0].b.ConnectClicked(func(checked bool) {
		buttonClicked(0)
	})
	AButtons[1].b.ConnectClicked(func(checked bool) {
		buttonClicked(1)
	})
	AButtons[2].b.ConnectClicked(func(checked bool) {
		buttonClicked(2)
	})
	AButtons[3].b.ConnectClicked(func(checked bool) {
		buttonClicked(3)
	})
	AButtons[4].b.ConnectClicked(func(checked bool) {
		buttonClicked(4)
	})
	AButtons[5].b.ConnectClicked(func(checked bool) {
		buttonClicked(5)
	})
	AButtons[6].b.ConnectClicked(func(checked bool) {
		buttonClicked(6)
	})
	AButtons[7].b.ConnectClicked(func(checked bool) {
		buttonClicked(7)
	})
	AButtons[8].b.ConnectClicked(func(checked bool) {
		buttonClicked(8)
	})

	slotsTabPage.SetLayout(slotsTabLayout)

	return slotsTabPage
}
