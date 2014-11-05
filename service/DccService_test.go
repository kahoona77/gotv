package service

import (
	"github.com/kahoona77/gotv/domain"
	"testing"
)

func TestGetTempFile(t *testing.T) {
	settings := new(domain.XtvSettings)
	settings.TempDir = "c:/temp"
	dcc := DccService{}
	fileEvent := DccFileEvent{"SEND", "simpsons.mkv", nil, "", 0}

	file := dcc.getTempFile(&fileEvent, settings)

	if file == nil {
		t.Error("no temp file")
	} else if file.Name() != "c:\\temp\\simpsons.mkv" {
		t.Error("wrong file: " + file.Name())
	}
}
