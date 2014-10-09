package handler

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/irc"
	"log"
	"net/http"
)

type DownloadsResult struct {
	Success  bool                `json:"success"`
	Status   string              `json:"status"`
	Packets  []domain.Packet     `json:"packets,omitempty"`
	Count    int                 `json:"count,omitempty"`
}

type DownloadsHandler struct {
	Client   *irc.IrcClient
}

func NewDownloadsHandler(client   *irc.IrcClient) *DownloadsHandler {
	h := new(DownloadsHandler)
	h.Client = client
	return h
}

func (this DownloadsHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.Print("URL: " + url)

	switch {
	case url== "/downloads/downloadPacket":
		this.downloadPacket(w, r)
	// case url == "/packets/listDownloads":
	// 	this.listDownloads(w, r)
	// case url == "/packets/cancelDownload":
	// 	this.cancelDownload(w, r)
	// case url == "/packets/resumeDownload":
	// 	this.resumeDownload(w, r)
	// case url == "/packets/stopDownload":
	// 	this.stopDownload(w, r)
	}

	return
}

func (this DownloadsHandler) downloadPacket(w http.ResponseWriter, r *http.Request) {
	var packet domain.Packet
	data := map[string]interface{}{
		"success": true,
		"status":  "ok",
	}
	if readJson(r, "data", &packet) {
		bot:= this.Client.GetBot (packet.Server)
		bot.DownloadPacket (&packet)
	}

	json.NewEncoder(w).Encode(data)
}
