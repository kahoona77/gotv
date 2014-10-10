package irc

import (
  "testing"
  "github.com/kahoona77/gotv/domain"
)

func TestGetTempFile(t *testing.T) {
  settings := new (domain.XtvSettings)
  settings.TempDir = "d:/temp"
  dcc := NewDccService (settings)
  fileEvent := DccFileEvent {"SEND", "simpsons.mkv", nil, "", 0}

  file := dcc.getTempFile (&fileEvent)

  if file == nil {
    t.Error("no temp file")
  } else if file.Name() != "d:\\temp\\simpsons.mkv" {
    t.Error("wrong file: " + file.Name())
  }
}
