package irc

import (
  "testing"
  "fmt"
)

func TestMakeOutBuffer(t *testing.T) {
  buf:=makeOutBuffer(234452)
  fmt.Printf("Buf: %v", buf)

  buf2:=makeOutBuffer(560980)
  fmt.Printf("Buf2: %v", buf2)

  if buf[0] != 0 {
    t.Error("buf error")
  }

  if buf2[0] != 0 {
    t.Error("buf error")
  }
}
