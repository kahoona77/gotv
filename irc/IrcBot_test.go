package irc

import (
	"fmt"
	"testing"
)

func TestMakeOutBuffer(t *testing.T) {
	buf := makeOutBuffer(234452)
	fmt.Printf("Buf: %v", buf)

	buf2 := makeOutBuffer(560980)
	fmt.Printf("Buf2: %v", buf2)

	if buf[0] != 0 || buf[1] != 3 || buf[2] != 147 || buf[3] != 212 {
		t.Error("buf error")
	}

	if buf2[0] != 0 || buf2[1] != 8 || buf2[2] != 143 || buf2[3] != 84 {
		t.Error("buf error")
	}
}
