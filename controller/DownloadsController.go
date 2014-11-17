package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/service"
)

type DownloadsResult struct {
	Success   bool                `json:"success"`
	Status    string              `json:"status"`
	Downloads []*service.Download `json:"downloads"`
}

type FilesResult struct {
	Success bool           `json:"success"`
	Status  string         `json:"status"`
	Files   []service.File `json:"files"`
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
	case url == "/downloads/loadFiles":
		dc.loadFiles(w, r)
	case url == "/downloads/deleteFiles":
		dc.deleteFiles(w, r)
	case url == "/downloads/moveFilesToMovies":
		dc.moveFilesToMovies(w, r)
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

func (dc DownloadsController) loadFiles(w http.ResponseWriter, r *http.Request) {

	files := dc.FilesService.GetFiles()
	data := FilesResult{true, "ok", files}
	json.NewEncoder(w).Encode(data)
}

func (dc DownloadsController) deleteFiles(w http.ResponseWriter, r *http.Request) {
	var files []service.File
	data := FilesResult{true, "ok", nil}
	if readJson(r, "data", &files) {
		err := dc.FilesService.DeleteFiles(files)
		if err != nil {
			log.Printf("ERROR: %v", err)
			data.Success = false
			data.Status = "Error while deleting files: " + err.Error()
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (dc DownloadsController) moveFilesToMovies(w http.ResponseWriter, r *http.Request) {
	var files []service.File
	data := FilesResult{true, "ok", nil}
	if readJson(r, "data", &files) {
		err := dc.FilesService.MoveFilesToMovies(files)
		if err != nil {
			log.Printf("ERROR: %v", err)
			data.Success = false
			data.Status = "Error while moving files: " + err.Error()
		}
	}

	json.NewEncoder(w).Encode(data)
}
