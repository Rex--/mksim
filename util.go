package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Two complement's add 2 x 12-bit unsigned integers stored as uint16's
// Returns a 12-bit usigned int stored as uint16, and a carry flag to signify an overflow has
// occurred.
func MKadd(a, b uint16) (x uint16, c bool) {
	c = false
	x = a + b
	// Check for 12-bit overflow
	if x > 4095 {
		c = true
		x = x - 4096
	}
	return x, c
}

func MKcomplement(a uint16) (x uint16) {

	// Complement the last 12 bits
	x = a ^ 0o7777

	// y := uint16(a)

	// for i := 0; i < 12; i++ {
	// 	x = (int16((1^a>>i)&1) << i) | x
	// }

	// // Detect signed int and convert it to negative
	// if ((x >> 11) & 1) == 1 {
	// 	// Keep bottom 11 bits
	// 	x &= 0b011111111111
	// 	// Set top bit to 1 (Convert to negative number)
	// 	x |= 0b100000000000
	// }

	return
}

func MKrotateRight(a uint16, l bool) (x uint16, y bool) {
	if a&1 == 1 {
		y = true
	}
	x = a >> 1
	if l {
		x |= 0b100000000000
	}
	return
}

func MKrotateLeft(a uint16, l bool) (x uint16, y bool) {
	if (a>>11)&1 == 1 {
		y = true
	}
	x = (a << 1) & 0b111111111111
	if l {
		x |= 0b1
	}
	return
}

// Load an object file produced by pdpnasm.
// This function returns an array of 4096 int16's representing pdp8 memory
func LoadPObjFile(filename string) (mem [4096]uint16, err error) {

	err = nil

	// Open file and create new scanner
	objFile, err := os.Open(filename)
	if err != nil {
		return
	}
	defer objFile.Close()
	scanner := bufio.NewScanner(objFile)

	// Loop over file line by line
	var addr uint16 = 0
	var data uint16
	var rawData uint64
	for scanner.Scan() {
		rawData, err = strconv.ParseUint(scanner.Text(), 8, 16)
		if err != nil {
			return
		}
		data = uint16(rawData)

		// If the 13th bit is set it's an address
		if data > 0o7777 {
			addr = (data & 0o7777)
		} else {
			mem[addr] = data
			addr++
		}
	}

	return
}

// Load a binary file in RIM format. These are produced by mkasm, but
// the RIM format was originally used for paper tapes for the PDP-8.
func LoadRIMFile(filename string) (mem [4096]uint16, err error) {

	// Open file and create a new reader
	rimFile, err := os.Open(filename)
	if err != nil {
		return
	}
	defer rimFile.Close()

	rimReader := bufio.NewReader(rimFile)

	// Skip over leading `0o200` bytes
	for b, _ := rimReader.Peek(1); b[0] == 0o200; b, _ = rimReader.Peek(1) {
		_, err = rimReader.Discard(1)
		if err != nil {
			return
		}
	}

	// Loop until EOF or trailing `0o200` bytes
	for {
		block := make([]byte, 4)
		read, err := rimReader.Read(block)
		if read != len(block) && read > 0 { // Sometimes err == EOF and read != 0
			if block[0] != 0o200 {
				return mem, err
				// panic("Incorrect format")
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return mem, err
			}
		}

		var addr, data uint16
		if b := block[0]; b == 0o200 { // Trailer bytes, break from loop
			break
		} else if b>>6&1 == 1 {
			// Start of address byte, this means the format is correct-ish
			addr = (uint16(block[0]&0o77) << 6) | uint16(block[1]&0o77)
			data = (uint16(block[2]&0o77) << 6) | uint16(block[3]&0o77)
		} else {
			return mem, fmt.Errorf("incorrect format")
			// panic("Incorrect format")
		}

		mem[addr] = data
	}

	return
}
