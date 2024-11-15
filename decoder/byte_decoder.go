package decoder

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	Reset = "\033[0m"
	Pink  = "\033[35m"
)

type Decoder struct {
	data   []byte
	offset int
	max    int
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		data:   data,
		offset: 0,
		max:    len(data),
	}
}

func LogMessage(data []byte, incoming bool) {
	// Log data as hex and ASCII
	x := hex.EncodeToString(data)
	xclean := ""
	for i, char := range x {
		if i%32 == 0 && i != 0 {
			xclean += "\n"
		} else if i%2 == 0 && i != 0 {
			xclean += " "
		}
		xclean += string(char)
	}
	if incoming {
		fmt.Printf("%sIncoming Data (hex):\n%s%s\n\n", Pink, xclean, Reset)
	} else {
		fmt.Printf("%sOutgoing Data (hex):\n%s%s\n\n", Pink, xclean, Reset)
	}
}

func (d *Decoder) readInt64() int64 {
	if d.offset+8 > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+8, d.max)
		return 0
	}
	n := buffToInt64(d.data[d.offset : d.offset+8])
	d.offset += 8
	return n
}

func (d *Decoder) readInt32() int32 {
	if d.offset+4 > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+4, d.max)
		return 0
	}
	n := buffToInt32(d.data[d.offset : d.offset+4])
	d.offset += 4
	return n
}

func (d *Decoder) readInt16() int16 {
	if d.offset+2 > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+2, d.max)
		return 0
	}
	n := buffToInt16(d.data[d.offset : d.offset+2])
	d.offset += 2
	return n
}

func (d *Decoder) readInt8() int8 {
	if d.offset+1 > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+1, d.max)
		return 0
	}
	n := buffToInt8(d.data[d.offset : d.offset+1])
	d.offset += 1
	return n
}

func (d *Decoder) readNullableString() string {
	length := d.readInt16()
	if length == -1 {
		return ""
	}
	if d.offset+int(length) > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+int(length), d.max)
		return ""
	}
	s := string(d.data[d.offset : d.offset+int(length)])
	d.offset += int(length)
	return s
}

func (d *Decoder) readCompactString() string {
	length := d.readInt8() - 1
	if length <= 0 {
		return ""
	}
	if d.offset+int(length) > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+int(length), d.max)
		return ""
	}
	s := string(d.data[d.offset : d.offset+int(length)])
	d.offset += int(length)
	return s
}

func buffToInt64(buff []byte) int64 {
	var n int64
	if err := binary.Read(bytes.NewReader(buff), binary.BigEndian, &n); err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}

func buffToInt32(buff []byte) int32 {
	var n int32
	if err := binary.Read(bytes.NewReader(buff), binary.BigEndian, &n); err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}

func buffToInt16(buff []byte) int16 {
	var n int16
	if err := binary.Read(bytes.NewReader(buff), binary.BigEndian, &n); err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}

func buffToInt8(buff []byte) int8 {
	var n int8
	if err := binary.Read(bytes.NewReader(buff), binary.BigEndian, &n); err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}
