package main

import (
	"errors"
	"go.bug.st/serial.v1"
	"go.bug.st/serial.v1/enumerator"
	"log"
	"strconv"
	"strings"
	"time"
)

var serialPort serial.Port
var SelectedPortId int
var SelectedDeviceId int

var SerialDevice1 string

func getSerialPorts() (usbports []string, perr error) {
	Devices.load()
	ports2, err := enumerator.GetDetailedPortsList()
	if err != nil {
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
			for di, d := range Devices.vendorId {
				if strings.ToLower(port.VID) == strings.ToLower(d) && strings.ToLower(port.PID) == strings.ToLower(Devices.productId[di]) {
					log.Printf("detected Device: %s\nportName: %s\n", Devices.cdc[di], port.Name)
					SelectedPortId = len(usbports) - 1
					SelectedDeviceId = di
					SerialDevice1 = port.Name
				}
			}
		}
	}

	return usbports, nil
}

func connectSerial(selSerialPort string) (err error) {
	c1 := make(chan int, 1)
	go func() {

		if selSerialPort == "" {
			err = errors.New("no device given")
		}

		mode := &serial.Mode{
			BaudRate: 115200,
			Parity:   serial.NoParity,
			DataBits: 8,
			StopBits: serial.OneStopBit,
		}
		serialPort, err = serial.Open(selSerialPort, mode)
		time.Sleep(time.Millisecond * 50)
		serialPort.ResetInputBuffer()
		serialPort.ResetOutputBuffer()
		if err != nil {
			log.Println("error serial connect ", err)
		} else {
			c1 <- 1
		}
	}()
	select {
	case res := <-c1:
		log.Printf("serialPort %v  connected - res: %d\n", serialPort, res)
	case <-time.After(time.Second * 2):
		err = errors.New("serial connection timeout")
	}

	if err != nil {
		return err
	}

	return
}

func sendSerialCmd(cmd string) {
	SerialResponse.Cmd = cmd
	SerialResponse.Code = -1
	SerialResponse.String = ""
	SerialResponse.Payload = ""

	temp := sendSerial(cmd)
	getSerialResponse(temp)
}

func sendSerial(cmdStr string) string {
	var resp string
	c1 := make(chan string)
	go func() {
		_, err := serialPort.Write([]byte(cmdStr + "\r\n"))
		if err != nil {
			log.Println("errro send serial: ", err, cmdStr)
		}
		time.Sleep(time.Millisecond * 25)
		resp = receiveSerial()
		c1 <- resp
	}()
	select {
	case resp := <-c1:
		return resp
	case <-time.After(time.Second * 5):
		log.Println("sendSrial Timeout")
	}
	return resp
}

func receiveSerial() (resp string) {
	buff := make([]byte, 512)
	var err error
	var n = 1
	var c = 0
	n, err = serialPort.Read(buff)
	c = c + n
	if err != nil {
		log.Printf("error temp: %s - n %d - error (%s)\n", resp, n, err)
	}
	return string(buff[:c])
}

func deviceInfo(longInfo string) (shortInfo string) {
	shortInfo = "undefined"
	toks := strings.Split(longInfo, " ")
	if len(toks) >= 3 {
		shortInfo = toks[0] + " " + toks[1]
	} else {
		shortInfo = longInfo
	}
	return
}

func getSerialResponse(res string) {
	var result []string

	res = strings.Replace(res, "\n", "#", -1)
	res = strings.Replace(res, "\r", "#", -1)
	res = strings.Replace(res, "##", "#", -1)
	if !strings.Contains(res, ":") {
		log.Println("no response given")
		Connected = false
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
	}
}
