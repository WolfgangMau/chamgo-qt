package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var configfile = "config.yaml"

type Config struct {
	Device []struct {
		Vendor  string `yaml:"vendor"`
		Product string `yaml:"product"`
		Name    string `yaml:"name"`
		Cdc     string `yaml:"cdc"`
		CmdSet 	map[string]string `yaml:"cmdset"`
		Config struct {
			Slot struct {
				Selected int `yaml:"selected"`
				Offset int `yaml:"offset"`
				First int `yaml:"first"`
				Last int `yaml:"last"`
			} `yaml:"slot"`
			Serial struct {
				Baud 				int `yaml:"baud"`
				WaitForReceive		int `yaml:"waitforreceive"`
				ConeectionTimeout 	int `yaml:"connectiontimeout"`
				Autoconnect 		bool `yaml:"autoconnect"`
			} `yaml:"serial"`
		} `yaml:"config"`
	} `yaml:"device"`
	Gui struct {
		Title string `yaml:"title"`
	} `yaml:"gui"`
}

func (c *Config) Load() *Config {
	cfgfile := Configpath()+configfile
	//log.Printf("Using configFile: %s\n", configfile)
	if len(cfgfile)>0 {
		yamlFile, err := ioutil.ReadFile(cfgfile)
		if err != nil {
			log.Printf("error reading config (%s) err   #%v ", cfgfile, err)
			os.Exit(2)
		}

		log.Println("loaded configfile: ",cfgfile)

		err = yaml.Unmarshal(yamlFile, c)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
			return nil
		}
		return c
	}
	return nil
}

func (c *Config) Save() bool {
	cfgfile := Configpath()+configfile
	if len(cfgfile)>0 {
		if data, err := yaml.Marshal(c); err != nil {
			log.Printf("error Marshall yaml (%s)\n", err)
			return false
		} else {
			err := ioutil.WriteFile(cfgfile, data, 0644)
			if err != nil {
				log.Fatal(err)
			}
			return true
		}
	}
	return false
}

func Configpath() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Printf("no executable-path found!?!\n%s\n",err)
		return ""
	}
	configpath := dir + string(filepath.Separator) + "config" + string(filepath.Separator)
	if _, err := os.Stat(configpath + configfile); os.IsNotExist(err) {
		log.Printf("ConfigFile %s not found!\n", configpath+configfile)
		return ""
	}
	return configpath
}

type DeviceActions struct {
	//config info
	GetModes    string
	GetButtons  string
	GetButtonsl string
	//slot info
	GetMode    string
	GetUid     string
	GetButton  string
	GetButtonl string
	GetSize    string
	//actions
	SelectSlot    string
	SelectedSlot  string
	ClearSlot     string
	StartUpload   string
	StartDownload string
}

func (d *DeviceActions) Load(commands map[string]string,device string) {
	switch device {

		case "Chameleon RevE-Rebooted":
			d.GetModes = commands["config"]
			d.GetButtons = commands["button"]

			d.GetMode = commands["config"] + "?"
			d.GetUid = commands["uid"] + "?"
			d.GetButton = commands["button"] + "?"
			d.GetButtonl = commands["buttonl"] + "?"
			d.GetSize = commands["memory"] + "?"

			d.SelectSlot = commands["setting"] + "="
			d.SelectedSlot = commands["setting"] + "?"
			d.StartUpload = commands["upload"]
			d.StartDownload = commands["download"]
			d.ClearSlot = commands["clear"]

		case "Chameleon RevG":
			d.GetModes = commands["config"] + "=?"
			d.GetButtons = commands["button"] + "=?"

			d.GetMode = commands["config"] + "?"
			d.GetUid = commands["uid"] + "?"
			d.GetButton = commands["button"] + "?"
			d.GetButtonl = commands["buttonl"] + "?"
			d.GetSize = commands["memory"] + "?"

			d.SelectSlot = commands["setting"] + "="
			d.SelectedSlot = commands["setting"] + "?"
			d.StartUpload = commands["upload"]
			d.StartDownload = commands["download"]
			d.ClearSlot = commands["clear"]
		}
}