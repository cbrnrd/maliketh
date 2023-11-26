package utils

import (
	"encoding/binary"
	"fmt"
	. "maliketh/pkg/config"
	"net"
	"strings"
	"time"
)

func DebugPrintln(msg string, args ...any) {
	if DEBUG {
		fmt.Printf(msg+"\n", args...)
	}
}

func ContainsAny(str string, elements []string) bool {
	for element := range elements {
		e := elements[element]
		if strings.Contains(str, e) {
			return true
		}
	}
	return false
}

func GetNTPTime() time.Time {
	type ntp struct {
		FirstByte, A, B, C uint8
		D, E, F            uint32
		G, H               uint64
		ReceiveTime        uint64
		J                  uint64
	}
	sock, err := net.Dial("udp", "us.pool.ntp.org:123")
	if err != nil {
		return time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	sock.SetDeadline(time.Now().Add((2 * time.Second)))
	defer sock.Close()
	transmit := new(ntp)
	transmit.FirstByte = 0x1b
	binary.Write(sock, binary.BigEndian, transmit)
	binary.Read(sock, binary.BigEndian, transmit)
	return time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(((transmit.ReceiveTime >> 32) * 1000000000)))
}
