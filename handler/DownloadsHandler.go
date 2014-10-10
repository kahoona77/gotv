package handler

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/irc"
	"log"
	"net/http"
	"strings"
)

type DownloadsResult struct {
	Success  bool                `json:"success"`
	Status   string              `json:"status"`
	Downloads  []*irc.Download    `json:"downloads"`
}

type DownloadsHandler struct {
	dcc   *irc.DccService
}

func NewDownloadsHandler(dcc   *irc.DccService) *DownloadsHandler {
	h := new(DownloadsHandler)
	h.dcc = dcc
	return h
}

func (this DownloadsHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.Print("URL: " + url)

	switch {
	case url== "/downloads/downloadPacket":
		this.downloadPacket(w, r)
	case strings.HasPrefix(url, "/downloads/listDownloads"):
		this.listDownloads(w, r)
	// case url == "/downloads/cancelDownload":
	// 	this.cancelDownload(w, r)
	// case url == "/downloads/resumeDownload":
	// 	this.resumeDownload(w, r)
	// case url == "/downloads/stopDownload":
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
		this.dcc.DownloadPacket (&packet)
	}

	json.NewEncoder(w).Encode(data)
}

func (this DownloadsHandler) listDownloads(w http.ResponseWriter, r *http.Request) {

	downloads:= this.dcc.ListDownloads()
	data := DownloadsResult {true,"ok", downloads}
	json.NewEncoder(w).Encode(data)
}
