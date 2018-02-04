package main

import (
	"errors"
	//newSerial "github.com/tarm/serial"
	"go.bug.st/serial.v1"
	"go.bug.st/serial.v1/enumerator"
	"log"
	//"reflect"
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
		//time.Sleep(time.Second * 1)

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
		//c := &newSerial.Config{Name: selSerialPort, Baud: 115200, ReadTimeout: time.Millisecond * 200}
		/*c := &newSerial.Config{Name: selSerialPort, Baud: 115200}
		serialPort, err = newSerial.OpenPort(c)*/
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
		//log.Println("serial connection timeout")
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
	var temp string
	c1 := make(chan string, 1)
	go func() {
		//time.Sleep(time.Second * 2)
		_, err := serialPort.Write([]byte(cmdStr + "\r\n"))
		if err != nil {
			log.Println("errro send serial: ", err, cmdStr)
		}
		//log.Printf("len = %d - crlf= %X\n", len([]byte("\r\n")), []byte("\r\n"))
		temp = receiveSerial()
		//log.Printf("serialReceive : %s\n", temp)
		//if len(temp) < 6 {
		//	temp = temp + receiveSerial()
		//}
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

func receiveSerial() (resp string) {
	buff := make([]byte, 1024)
	var err error
	var n = 1
	var c = 0

	for n > 0 {
		c++
		n, _ = serialPort.Read(buff)
		//log.Printf("n: %d\n", n)
		if err != nil {
			log.Printf("error temp: %s - n %d - error (%s)\n", resp, n, err)
		}
		//log.Printf("n: %d\n", n)
		if n > 0 {
			resp = resp + string(buff[:n])
			//check if there is a CRLF at the end
			if len(resp) >= 2 && resp[(len(resp)-2):] == string([]byte{0x0D, 0x0A}) {
				//log.Println("receiveSerial -> EOL")
				n = 0
			}
			//log.Printf("fresh read of %d bytes - temp: %s\n", n, resp)
		}
	}

	//log.Printf("end of loop %d -  n: %d - Buffer: %s\n", c, n, resp)
	n = len(resp)
	//log.Printf("loops= %d  bytes= %d\n", c, n)
	return resp
}

func deviceInfo(longInfo string) (shortInfo string) {
	shortInfo = "undefined"
	toks := strings.Split(longInfo, " ")
	if len(toks) >= 3 {
		shortInfo = toks[0] + " " + toks[1]
	} else {
		shortInfo = longInfo
	}

	//log.Printf("Long : %s -> toks : %d -> short : %s\n", longInfo, len(toks), shortInfo)
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
		//dumpResp(SerialResponse)
	}
}

//func dumpResp(t serialResponse) {
//	s := reflect.ValueOf(&t).Elem()
//	typeOfT := s.Type()
//
//	for i := 0; i < s.NumField(); i++ {
//		f := s.Field(i)
//		log.Printf("%d: %s %s = %v\n", i,
//			typeOfT.Field(i).Name, f.Type(), f.Interface())
//	}
//}
