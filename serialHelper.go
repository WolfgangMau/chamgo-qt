package main

import (
	"go.bug.st/serial.v1"
	"log"
	"strings"
	"strconv"
	"reflect"
)

var serialPort serial.Port

func getSerialPorts() (ports []string, perr error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if len(ports) == 0 {
		log.Println("No serial ports found!")
		return nil, nil
	}

	return ports, nil
}

func connectSerial(selSerialPort string) (err error) {
	if selSerialPort == "" {
		return nil
	}
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	serialPort, err = serial.Open(selSerialPort, mode)

	if err != nil {
		return err
	}
	return nil
}

func sendSerialCmd(cmd string) {
	//reset response-struct
	SerialResponse.Cmd=cmd
	SerialResponse.Code=-1
	SerialResponse.String=""
	SerialResponse.Payload=""

	temp := sendSerial(cmd)
	getSerialResponse(temp)
}

func sendSerial(cmdStr string) string {
	//myLog(cmdStr + "\n")
	_, err := serialPort.Write([]byte(cmdStr + "\r"))
	if err != nil {
		log.Fatal(err)
	}
	res := receiveSerial()
	if res == "1" {
		res = "1" + receiveSerial()
	}
	return res
}

func receiveSerial() (recv string) {
	buff := make([]byte, 1024)
	for {
		// Reads up to 1024 bytes
		n, err := serialPort.Read(buff)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		if n <= 0 {
			log.Println("\nEOF")
			return ""
		}
		return string(buff[:n])
	}
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
	res = strings.Replace(res, "\n","#",-1)
	res = strings.Replace(res, "\r","#",-1)
	res = strings.Replace(res,"##","#",-1)

	temp2 = strings.Split(res, ":")
	if len(temp2[1]) >= 2 {
		result = append(result, temp2[0])
		SerialResponse.Code,_ = strconv.Atoi(temp2[0])
		temp := strings.Split(temp2[1], "#")
		if len(temp) > 0 {
			for i, s := range temp {
				switch i {
				case 0 :
					SerialResponse.String = s
				case 1 :
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