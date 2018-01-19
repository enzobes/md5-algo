package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	password = kingpin.Flag("password", "Password to hash.").Short('p').Required().String()
)

func main() {
	// using extrenal library for flag parsing
	kingpin.Parse()

	// hashing
	hash := md5(*password)

	//printing
	fmt.Printf("Hash de \"%s\": ", *password)
	for i := 0; i < len(hash); i++ {
		fmt.Printf("%x", hash[i])
	}
}

// shifting consts
var shift = [...]uint{7, 12, 17, 22, 5, 9, 14, 20, 4, 11, 16, 23, 6, 10, 15, 21}
var table [64]uint32

// constants definition from part of the sines of integers
func init() {
	for i := range table {
		table[i] = uint32((1 << 32) * math.Abs(math.Sin(float64(i+1))))
	}
}

func md5(s string) (r [16]byte) {
	padded := bytes.NewBuffer([]byte(s))
	padded.WriteByte(0x80)
	for padded.Len()%64 != 56 {
		padded.WriteByte(0)
	}
	messageLenBits := uint64(len(s)) * 8
	binary.Write(padded, binary.LittleEndian, messageLenBits)

	// main variables
	var a, b, c, d uint32 = 0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476

	var buffer [16]uint32

	// proccessing of the message
	for binary.Read(padded, binary.LittleEndian, buffer[:]) == nil { // read every 64 bytes

		// initialize hash values for this chunk
		a1, b1, c1, d1 := a, b, c, d
		// main loop
		for j := 0; j < 64; j++ {
			var f uint32
			bufferIndex := j
			round := j >> 4
			switch round {
			case 0:
				f = (b1 & c1) | (^b1 & d1)
			case 1:
				f = (b1 & d1) | (c1 & ^d1)
				bufferIndex = (bufferIndex*5 + 1) & 0x0F
			case 2:
				f = b1 ^ c1 ^ d1
				bufferIndex = (bufferIndex*3 + 5) & 0x0F
			case 3:
				f = c1 ^ (b1 | ^d1)
				bufferIndex = (bufferIndex * 7) & 0x0F
			}
			// rotations
			sa := shift[(round<<2)|(j&3)]
			a1 += f + buffer[bufferIndex] + table[j]
			a1, d1, c1, b1 = d1, c1, b1, a1<<sa|a1>>(32-sa)+b1
		}
		// adding chunk hash to result
		a, b, c, d = a+a1, b+b1, c+c1, d+d1
	}

	// little-endian output
	binary.Write(bytes.NewBuffer(r[:0]), binary.LittleEndian, []uint32{a, b, c, d})
	return
}
