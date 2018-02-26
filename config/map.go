package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

const MYFAIRE_CLASSIC_1K_4Byte_UID = 0x0400

type TagMap struct {
	Name string `yaml:"name"`
	Tagtype string `yaml:"tagtype"`
	Mappings []ByteMap `yaml:"mappings"`
}
type ByteMap struct {
	Start int `yaml:"start"`
	End int `yaml:"end"`
	MapBytes []MapByte `yaml:"mapbytes"`
	MapFuncs []MapFunc `yaml:"mapfuncs"`
}
type MapByte struct {
	Pos int `yaml:"pos"`
	Color string `yaml:"color"`
	Tooltip string `yaml:"tooltip"`
}
type MapFunc struct{
	Name string `yaml:"name"`
	ExpectResult bool `yaml:"expectresult"`
	Result []string `yaml:"result"`
}

var DefaultMap = TagMap{
	Name: "Mifare Classic 1K 4Byte UID",
	Tagtype: "0400",
	Mappings: []ByteMap{
		{
			Start: 0,
			End:3,
			MapBytes: []MapByte{
				{
					Pos: 0,
					Color: "yellow",
					Tooltip: "UID0",
				},
				{
					Pos: 1,
					Color: "yellow",
					Tooltip: "UID1",
				},
				{
					Pos: 2,
					Color: "yellow",
					Tooltip: "UID2",
				},
				{
					Pos: 3,
					Color: "yellow",
					Tooltip: "UID3",
				},
			},
		},
		{
			Start: 4,
			End:4,
			MapBytes: []MapByte{
				{
					Pos: 0,
					Color: "yellow",
					Tooltip: "BCC (UID0..UID3)",
				},
			},
			MapFuncs: []MapFunc{
				{
					Name: "bcc",
					ExpectResult: true,
					Result: []string{"fc"},
				},
			},
		},
	},
}

func (m *TagMap) Save(f string) {
	tagmap := Apppath()+string(filepath.Separator)+runtime.GOOS+string(filepath.Separator)+"maps"+string(filepath.Separator)+f
	if len(m.Mappings) > 0 {
		if data, err := yaml.Marshal(m); err != nil {
			log.Printf("error Marshall yaml (%s)\n", err)
		} else {
			err := ioutil.WriteFile(tagmap, data, 0644)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func (c *TagMap) Load(f string) {
	tagmap := Apppath()+string(filepath.Separator)+runtime.GOOS+string(filepath.Separator)+"maps"+string(filepath.Separator)+f
	//log.Printf("Using configFile: %s\n", configfile)
	if len(tagmap) > 0 {
		yamlFile, err := ioutil.ReadFile(tagmap)
		if err != nil {
			log.Printf("error reading config (%s) err   #%v ", tagmap, err)
			return
		}
		log.Println("loaded TagMap: ", c.Name)

		err = yaml.Unmarshal(yamlFile, c)
		if err != nil {
			log.Print("Unmarshal: %v", err)
		}
	}
}