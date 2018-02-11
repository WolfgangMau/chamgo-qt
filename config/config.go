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

//func  (c *Config) Default() *Config {
//	//Device
//	return
//}

func (c *Config) Load() *Config {
	cfgfile := configpath()
	//log.Printf("Using configFile: %s\n", configfile)
	if len(cfgfile)>0 {
		yamlFile, err := ioutil.ReadFile(cfgfile)
		if err != nil {
			log.Printf("error reading config (%s) err   #%v ", cfgfile, err)
			os.Exit(2)
		}
		//log.Printf("%v\n",string(yamlFile))

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
	cfgfile := configpath()
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

func configpath() string {

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
	return configpath+configfile
}