package nonces

import "log"

type Nonce struct {
	Key    byte
	Sector byte
	Nt     []byte
	Nr     []byte
	Ar     []byte
}

func DecryptData(encarr []byte, key int, size int) []byte {
	arr := make([]byte, size)
	arr = encarr
	for i := 0; i < size; i++ {
		s := int(arr[i])
		t := size + key + i - size/key ^ s
		encarr[i] = byte(t)
	}
	return encarr
}

func ExtractNonces(data []byte) (res []Nonce) {
	for i := 16; i < (208 - 16); i = i + 16 {
		var n Nonce
		n.Key = data[i]          //16
		n.Sector = data[i+1]     //17
		n.Nt = data[i+4 : i+8]   //20-23
		n.Nr = data[i+8 : i+12]  //24-27
		n.Ar = data[i+12 : i+16] //28-31
		if n.Key != byte(0xff) && n.Sector != byte(0xff) {
			res = append(res, n)
			log.Printf("key: %x  sector: %x %x %x %x\n",n.Key, n.Sector, n.Nt, n.Nr, n.Ar)
		}
	}
	return res
}
