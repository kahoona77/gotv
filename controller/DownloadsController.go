package controller

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/service"
	"net/http"
	"strings"
)

type DownloadsResult struct {
	Success   bool                `json:"success"`
	Status    string              `json:"status"`
	Downloads []*service.Download `json:"downloads"`
}

type DownloadsController struct {
	*service.Context
}


func (dc DownloadsController) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	switch {
	case url == "/downloads/downloadPacket":
		dc.downloadPacket(w, r)
	case strings.HasPrefix(url, "/downloads/listDownloads"):
		dc.listDownloads(w, r)
	case url == "/downloads/cancelDownload":
		dc.cancelDownload(w, r)
	case url == "/downloads/resumeDownload":
		dc.resumeDownload(w, r)
	case url == "/downloads/stopDownload":
		dc.stopDownload(w, r)
	}

	return
}

func (dc DownloadsController) downloadPacket(w http.ResponseWriter, r *http.Request) {
	var packet domain.Packet
	data := map[string]interface{}{
		"success": true,
		"status":  "ok",
	}
	if readJson(r, "data", &packet) {
		dc.DccService.DownloadPacket(&packet)
	}

	json.NewEncoder(w).Encode(data)
}

func (dc DownloadsController) listDownloads(w http.ResponseWriter, r *http.Request) {

	downloads := dc.DccService.ListDownloads()
	data := DownloadsResult{true, "ok", downloads}
	json.NewEncoder(w).Encode(data)
}

func (dc DownloadsController) stopDownload(w http.ResponseWriter, r *http.Request) {
	var download service.Download
	if readJson(r, "data", &download) {
		dc.DccService.StopDownload(&download)
	}

	data := DownloadsResult{true, "ok", nil}
	json.NewEncoder(w).Encode(data)
}

func (dc DownloadsController) cancelDownload(w http.ResponseWriter, r *http.Request) {
	var download service.Download
	if readJson(r, "data", &download) {
		dc.DccService.CancelDownload(&download)
	}

	data := DownloadsResult{true, "ok", nil}
	json.NewEncoder(w).Encode(data)
}

func (dc DownloadsController) resumeDownload(w http.ResponseWriter, r *http.Request) {
	var download service.Download
	if readJson(r, "data", &download) {
		dc.DccService.ResumeDownload(&download)
	}

	data := DownloadsResult{true, "ok", nil}
	json.NewEncoder(w).Encode(data)
}
