package main

import (
	"fmt"
	"go.bug.st/serial.v1"
	"log"
	"strings"
)

var serialPort serial.Port

func getSerialPorts() (ports []string, perr error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
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

func serialCMD(cmd string) []string {
	log.Print("-> " + cmd)
	temp := sendSerial(cmd)
	return responseSplit(temp)
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
			fmt.Println("\nEOF")
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

func responseSplit(res string) []string {
	var result []string
	res = strings.Replace(res, "\n","#",-1)
	res = strings.Replace(res, "\r","#",-1)
	res = strings.Replace(res,"##","#",-1)

	temp2 = strings.Split(res, ":")
	if len(temp2[1]) >= 2 {
		result = append(result, temp2[0])
		temp := strings.Split(temp2[1], "#")
		if len(temp) > 0 {
			for _, s := range temp {
				if s != "" {
					result = append(result, s)
				}
			}
		}
		for i,s := range result {
			log.Printf("result[%d] = %s\n",i,s)
		}
		return result
	} else {
		log.Printf("result has no reslutcode! %s\n",res)
	}
	return nil
}