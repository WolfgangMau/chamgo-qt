// Package crc16 implements the 16-bit cyclic redundancy check, or CRC-16, checksum.
//
// It provides parameters for the majority of well-known CRC-16 algorithms.
package crc16

import (
	"encoding/hex"
	"fmt"
)

// Params represents parameters of CRC-16 algorithms.
// More information about algorithms parametrization and parameter descriptions
// can be found here - http://www.zlib.net/crc_v3.txt
type Params struct {
	Poly   uint16
	Init   uint16
	RefIn  bool
	RefOut bool
	XorOut uint16
	Check  uint16
	Name   string
}

// Predefined CRC-16 algorithms.
var (
	CRC16_CRC_A = Params{0x1021, 0xC6C6, true, true, 0x0000, 0xBF05, "CRC-16/CRC-A"}
)

// Table is a 256-word table representing polinomial and algorithm settings for efficient processing.
type Table struct {
	params Params
	data   [256]uint16
}

// MakeTable returns the Table constructed from the specified algorithm.
func MakeTable(params Params) *Table {
	table := new(Table)
	table.params = params
	for n := 0; n < 256; n++ {
		crc := uint16(n) << 8
		for i := 0; i < 8; i++ {
			bit := (crc & 0x8000) != 0
			crc <<= 1
			if bit {
				crc ^= params.Poly
			}
		}
		table.data[n] = crc
	}
	return table
}

// Init returns the initial value for CRC register corresponding to the specified algorithm.
func Init(table *Table) uint16 {
	return table.params.Init
}

// Update returns the result of adding the bytes in data to the crc.
func Update(crc uint16, data []byte, table *Table) uint16 {
	for _, d := range data {
		if table.params.RefIn {
			d = ReverseByte(d)
		}
		crc = crc<<8 ^ table.data[byte(crc>>8)^d]
	}
	return crc
}

// Complete returns the result of CRC calculation and post-calculation processing of the crc.
func Complete(crc uint16, table *Table) uint16 {
	if table.params.RefOut {
		return ReverseUint16(crc) ^ table.params.XorOut
	}
	return crc ^ table.params.XorOut
}

// Checksum returns CRC checksum of data usign scpecified algorithm represented by the Table.
func Checksum(data []byte, table *Table) uint16 {
	crc := Init(table)
	crc = Update(crc, data, table)
	return Complete(crc, table)
}

func ReverseByte(val byte) byte {
	var rval byte = 0
	for i := uint(0); i < 8; i++ {
		if val&(1<<i) != 0 {
			rval |= 0x80 >> i
		}
	}
	return rval
}

func ReverseUint16(val uint16) uint16 {
	var rval uint16 = 0
	for i := uint(0); i < 16; i++ {
		if val&(uint16(1)<<i) != 0 {
			rval |= uint16(0x8000) >> i
		}
	}
	return rval
}

//calculate BCC of a Mifare Classic
func GetBCC(input []byte) (bcc byte) {
	bcc = byte(0x0)

	if len(input) > 0 {
		for i := 0; i < len(input); i++ {
			bcc ^= input[i]
		}
	}
	return bcc
}

func Split2Hex(s string) []byte {
	res, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return res
}

func GetCRCA(hexstr string) string {
	strhex := Split2Hex(hexstr)
	t := MakeTable(CRC16_CRC_A)
	checksum5 := Checksum(strhex, t)
	var h, l uint8 = uint8(checksum5 >> 8), uint8(checksum5 & 0xff)
	return fmt.Sprintf("%02X%02X", l, h)
}
