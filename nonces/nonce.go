package nonces

type Nonce struct {
	key    byte
	sector byte
	nt     []byte
	nr     []byte
	ar     []byte
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
		n.key = data[i]          //16
		n.sector = data[i+1]     //17
		n.nt = data[i+4 : i+8]   //20-23
		n.nr = data[i+8 : i+12]  //24-27
		n.ar = data[i+12 : i+16] //28-31
		if n.key != byte(0xff) && n.sector != byte(0xff) {
			res = append(res, n)
		}
	}
	return res
}
