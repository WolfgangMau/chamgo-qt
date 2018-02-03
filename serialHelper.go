package main

import (
	"go.bug.st/serial.v1"
	newSerial "github.com/tarm/serial"
	"go.bug.st/serial.v1/enumerator"
	"log"
	"strconv"
	"strings"
	"time"
	"errors"
	"reflect"
	"fmt"
)

var serialPort *newSerial.Port
var SelectedPortId int
var SelectedDeviceId int

var SerialDevice1 string

func getSerialPorts() (usbports[]string, perr error) {
	Devices.load()
	ports2,err := enumerator.GetDetailedPortsList()
	if err != nil  {
		log.Fatal(err)
		return nil, err
	}
	ports, err2 := serial.GetPortsList()
	if err2 != nil {
		log.Fatal(err)
		return nil, err2
	}
	if err != nil || err2 != nil {
		log.Fatal(err)
		return nil, err
	}
	if len(ports) == 0 {
		log.Println("No serial ports found!")
		return nil, nil
	}

	for _, port := range ports2 {
		if port.IsUSB {
			usbports = append(usbports, port.Name)
			for di,d  := range Devices.vendorId {
				if strings.ToLower(port.VID) == strings.ToLower(d) && strings.ToLower(port.PID) == strings.ToLower(Devices.productId[di]) {
					log.Printf("detected Device: %s\nportName: %s\n",Devices.cdc[di], port.Name)
					SelectedPortId = len(usbports)-1
					SelectedDeviceId = di
					SerialDevice1 = port.Name
				}
			}
		}
	}

	return usbports, nil
}

func connectSerial(selSerialPort string) ( err error) {
	c1 := make(chan int, 1)
	go func() {
		//time.Sleep(time.Second * 1)

		if selSerialPort == "" {
			err = errors.New("no device given")
		}

		//mode := &serial.Mode{
		//	BaudRate: 115200,
		//	Parity:   serial.NoParity,
		//	DataBits: 8,
		//	StopBits: serial.OneStopBit,
		//}
		//serialPort, err = serial.Open(selSerialPort, mode)
		c := &newSerial.Config{Name: selSerialPort, Baud: 115200}
		serialPort,err = newSerial.OpenPort(c)
		if err != nil {
			log.Println("error serial connect ",err)
		} else {
			c1 <- 1
		}
	}()
	select {
	case res := <-c1:
		log.Printf("serialPort %v  connected - res: %d\n", serialPort, res)
	case <-time.After(time.Second * 2):
		log.Println("serial connection timeout")
		err = errors.New("serial connection timeout")
	}

	if err != nil {
		return err
	}

	return
}

func sendSerialCmd(cmd string) {
	//reset response-struct
	SerialResponse.Cmd = cmd
	SerialResponse.Code = -1
	SerialResponse.String = ""
	SerialResponse.Payload = ""

	temp := sendSerial(cmd)
	getSerialResponse(temp)
}

func sendSerial(cmdStr string) string {
	serialPort.Flush()
	var temp string
	c1 := make(chan string, 1)
	go func() {
		//time.Sleep(time.Second * 2)
		_, err := serialPort.Write([]byte(cmdStr + "\r\n"))
		if err != nil {
			log.Println("errro send serial: ",err, cmdStr)
		}
		//log.Print("n: ",n)
		temp = receiveSerial()
		if temp == "1" {
			temp = "1" + receiveSerial()
		}
		c1 <- temp
	}()
	select {
	case res := <-c1:
		return res
	case <-time.After(time.Second * 5):
		log.Println("sendSrial Timeout")
	}
	return temp
}

func receiveSerial() (recv string) {
	buff := make([]byte, 512)
	n:=0
	var err error
	//var err error
	for {
		log.Println("loop receive")
		// Reads up to 512 bytes
		n, err = serialPort.Read(buff)
		if err != nil {
			log.Println("errro receive serial: ",err)
		}
		log.Printf("%q", buff[:n])
		//minimum 6 bytes reqierd : 101:OK
		if n > 5 {
			break
		}
	}
	log.Println("receive: ",fmt.Sprintf("%s",buff[:n]),n)
	return string(buff[:n])
}

func deviceInfo(longInfo string) (shortInfo string) {
	shortInfo = "undefined"
	toks := strings.Split(longInfo, " ")
	if len(toks) >= 3 {
		shortInfo = toks[0] + " " + toks[1]
	} else {
		shortInfo = longInfo
	}

	log.Printf("Long : %s -> toks : %d -> short : %s\n", longInfo, len(toks), shortInfo)
	return
}

func getSerialResponse(res string) {
	var result []string
	log.Printf("q getSerialResponse: %q\n",res)
	log.Printf("s getSerialResponse: %s\n",res)

	res = strings.Replace(res, "\n", "#", -1)
	res = strings.Replace(res, "\r", "#", -1)
	res = strings.Replace(res, "##", "#", -1)
	if !strings.Contains(res,":") {
		log.Println("no response given")
		Connected=false
		return
	}
	temp2 = strings.Split(res, ":")
	if len(temp2[1]) >= 2 {
		result = append(result, temp2[0])
		SerialResponse.Code, _ = strconv.Atoi(temp2[0])
		temp := strings.Split(temp2[1], "#")
		if len(temp) > 0 {
			for i, s := range temp {
				switch i {
				case 0:
					SerialResponse.String = s
				case 1:
					SerialResponse.Payload = s
				}
				if s != "" {
					result = append(result, s)
				}
			}
		}
		dumpResp(SerialResponse)
	}
}

func dumpResp(t serialResponse) {
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		log.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}
