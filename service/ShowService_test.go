package service

import (
	"testing"
)

func TestParseShow(t *testing.T) {
	parser := NewShowService(nil)
	info := parser.parseShow("d:/test/downloads/Breaking.Bad.S05E15.HDTV.x264.mkv")
	checkInfo(info, t)

	info = parser.parseShow("d:/test/downloads/Breaking.Bad.5x15.HDTV.x264.mkv")
	checkInfo(info, t)
}

func checkInfo(info *ShowInfo, t *testing.T) {
	if info == nil {
		t.Error("info is nil")
		return
	}

	if info.Name != "Breaking Bad" {
		t.Errorf("wrong name: %v", info.Name)
	}

	if info.Season != 5 {
		t.Errorf("wrong season: %v", info.Season)
	}

	if info.Episode != 15 {
		t.Errorf("wrong Episode: %v", info.Episode)
	}
}

func TestCleanFileName(t *testing.T) {
	filename := "Modern Family - 5x05 - Test: This is a test.mkv"
	result := sanitizeFilename(filename)

	if result != "Modern Family - 5x05 - Test This is a test.mkv" {
		t.Errorf("file not cleaned '%v'", result)
	}
}
