package main

var populated = false
var TagModes []string
var TagButtons []string

var ActionButtons = []string{"Select All", "Select None", "Apply", "Clear", "Refresh", "Set Active", "mfkey32", "Upload", "Download"}

var SerialResponse serialResponse

type serialResponse struct {
	Cmd     string
	Code    int
	String  string
	Payload string
}

var DeviceActions deviceActions

type deviceActions struct {
	//config info
	getModes    string
	getButtons  string
	getButtonsl string
	//slot info
	getMode    string
	getUid     string
	getButton  string
	getButtonl string
	getSize    string
	//actions
	selectSlot    string
	selectedSlot  string
	clearSlot     string
	startUpload   string
	startDownload string
}

func (d *deviceActions) load(device string) {
	Commands := Cfg.Device[SelectedDeviceId].CmdSet
	switch device {

	case "Chameleon RevE-Rebooted":
		d.getModes = Commands["config"]
		d.getButtons = Commands["button"]

		d.getMode = Commands["config"] + "?"
		d.getUid = Commands["uid"] + "?"
		d.getButton = Commands["button"] + "?"
		d.getButtonl = Commands["buttonl"] + "?"
		d.getSize = Commands["memory"] + "?"

		d.selectSlot = Commands["setting"] + "="
		d.selectedSlot = Commands["setting"] + "?"

		d.startUpload = Commands["upload"]
		d.startDownload = Commands["download"]
		d.clearSlot = Commands["clear"]

	case "Chameleon RevG":
		d.getModes = Commands["config"] + "=?"
		d.getButtons = Commands["button"] + "=?"

		d.getMode = Commands["config"] + "?"
		d.getUid = Commands["uid"] + "?"
		d.getButton = Commands["button"] + "?"
		d.getButtonl = Commands["buttonl"] + "?"
		d.getSize = Commands["memory"] + "?"

		d.selectSlot = Commands["setting"] + "="
		d.selectedSlot = Commands["setting"] + "?"
		d.startUpload = Commands["upload"]
		d.startDownload = Commands["download"]
		d.clearSlot = Commands["clear"]
	}
}
