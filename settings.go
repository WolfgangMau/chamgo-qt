package main

var populated = false
var TagModes []string
var TagButtons []string

//var Device string
//var Devices devices
//
//type devices struct {
//	name      []string
//	vendorId  []string
//	productId []string
//	cdc       []string
//}

////ToDo: use map map[string][string}
//func (d *devices) load() {
//	d.name = []string{"RevE-Rebooted", "RevG"}
//	d.vendorId = []string{"03eb", "16d0"}
//	d.productId = []string{"2044", "04b2"}
//	d.cdc = []string{"Atmel Corporation LUFA CDC", "KAOS Chameleon-Mini"}
//}

var ActionButtons = []string{"Select All", "Select None", "Apply", "Clear", "Refresh", "Set Active", "mfkey32", "Upload", "Download"}

var SerialResponse serialResponse

type serialResponse struct {
	Cmd     string
	Code    int
	String  string
	Payload string
}

//var Commands commands
//
//type commands struct {
//	version       string
//	config        string
//	uid           string
//	readonly      string
//	upload        string
//	download      string
//	reset         string
//	upgrade       string
//	memory        string
//	uidsize       string
//	button        string
//	buttonl       string
//	rbutton       string
//	rbuttonl      string
//	lbutton       string
//	lbuttonl      string
//	setting       string
//	clear         string
//	help          string
//	rssi          string
//	ledgreen      string
//	ledred        string
//	logmode       string
//	logmem        string
//	logdownload   string
//	logstore      string
//	logclear      string
//	store         string
//	recall        string
//	charging      string
//	systick       string
//	sendraw       string
//	send          string
//	getuid        string
//	dumpmfu       string
//	identify      string
//	timeout       string
//	threshold     string
//	autocalibrate string
//	field         string
//}
//
//func (c *commands) load(device string) {
//	switch device {
//	case "RevE-Rebooted":
//		c.version = "VERSIONMY"
//		c.config = "CONFIGMY"
//		c.uid = "UIDMY"
//		c.readonly = "READONLYMY"
//		c.upload = "UPLOADMY"
//		c.download = "DOWNLOADMY"
//		c.reset = "RESETMY"
//		c.upgrade = "UPGRADEMY"
//		c.memory = "MEMSIZEMY"
//		c.uidsize = "UIDSIZEMY"
//		c.button = "BUTTONMY"
//		c.buttonl = "BUTTON_LONGMY"
//		c.setting = "SETTINGMY"
//		c.clear = "CLEARMY"
//		c.help = "HELPMY"
//		c.rssi = "RSSIMY"
//
//	//http://rawgit.com/emsec/ChameleonMini/master/Doc/Doxygen/html/_page__command_line.html
//	case "RevG":
//		c.version = "VERSION"
//		c.config = "CONFIG"
//		c.uid = "UID"
//		c.readonly = "READONLY"
//		c.upload = "UPLOAD"
//		c.download = "DOWNLOAD"
//		c.reset = "RESET"
//		c.upgrade = "UPGRADE"
//		c.memory = "MEMSIZE"
//		c.uidsize = "UIDSIZE"
//		c.button = "RBUTTON"
//		c.buttonl = "RBUTTON_LONG"
//		c.lbutton = "LBUTTON"
//		c.lbuttonl = "LBUTTON_LONG"
//		c.setting = "SETTING"
//		c.clear = "CLEAR"
//		c.help = "HELP"
//		c.rssi = "RSSI"
//		c.ledgreen = "LEDGREEN"
//		c.ledred = "LEDRED"
//		c.logmode = "LOGMODE"
//		c.logmem = "LOGMEM"
//		c.logdownload = "LOGDOWNLOAD"
//		c.logstore = "LOGSTORE"
//		c.logclear = "LOGCLEAR"
//		c.store = "STORE"
//		c.recall = "RECALL"
//		c.charging = "CHARGING"
//		c.systick = "SYSTICK"
//		c.sendraw = "SEND_RAW"
//		c.send = "SEND"
//		c.getuid = "GETUID"
//		c.dumpmfu = "DUMP_MFU"
//		c.identify = "IDENTIFY"
//		c.timeout = "TIMEOUT"
//		c.threshold = "THRESHOLD"
//		c.autocalibrate = "AUTOCALIBRATE"
//		c.field = "FIELD"
//	}
//}

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

	case "RevE-Rebooted":
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

	case "RevG":
		d.getModes = Commands["config"] + "=?"
		d.getButtons = Commands["button"] + "=?"

		d.getMode = Commands["config"] + "?"
		d.getUid = Commands["uid"] + "?"
		d.getButton = Commands["button"] + "?"
		d.getButtonl = Commands["button"] + "?"
		d.getSize = Commands["memory"] + "?"

		d.selectSlot = Commands["setting"] + "="
		d.selectedSlot = Commands["setting"] + "?"
		d.startUpload = Commands["upload"]
		d.startDownload = Commands["download"]
		d.clearSlot = Commands["clear"]
	}
}
