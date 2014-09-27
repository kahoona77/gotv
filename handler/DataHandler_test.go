package handler

import (
	"github.com/kahoona77/gotv/domain"
	"testing"
)

func TestRemoveFromChannel(t *testing.T) {
	channels := []domain.Channel{domain.Channel{Name: "test1"}, domain.Channel{Name: "test2"}, domain.Channel{Name: "test3"}}
	removeChannel(&channels, channels[1])

	if len(channels) != 2 {
		t.Error("remove did not work")
	}
}
