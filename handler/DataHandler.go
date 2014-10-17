package handler

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/irc"
	"log"
	"net/http"
)

type DataResult struct {
	Success  bool                `json:"success"`
	Status   string              `json:"status"`
	Servers  []domain.Server     `json:"servers,omitempty"`
	Settings *domain.XtvSettings `json:"settings,omitempty"`
}

type DataHandler struct {
	ServerRepo   *domain.GoTvRepository
	SettingsRepo *domain.GoTvRepository
	DccService   *irc.DccService
}

func NewDataHandler(serverRepo *domain.GoTvRepository, settingsRepo *domain.GoTvRepository, dccService *irc.DccService) *DataHandler {
	h := new(DataHandler)
	h.ServerRepo = serverRepo
	h.SettingsRepo = settingsRepo
	h.DccService = dccService
	return h
}

func (this DataHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.Print("URL: " + url)

	switch {
	case url == "/data/loadServers":
		this.loadServers(w, r)
	case url == "/data/saveServer":
		this.saveServer(w, r)
	case url == "/data/deleteServer":
		this.deleteServer(w, r)
	case url == "/data/loadSettings":
		this.loadSettings(w, r)
	case url == "/data/saveSettings":
		this.saveSettings(w, r)
	}

	return
}

func (this DataHandler) saveServer(w http.ResponseWriter, r *http.Request) {
	var (
		server domain.Server
		err    error
	)
	data := DataResult{true, "ok", nil, nil}
	if readJson(r, "data", &server) {

		_, err = this.ServerRepo.Save(server.Id, &server)

		if err != nil {
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (this DataHandler) deleteServer(w http.ResponseWriter, r *http.Request) {
	var (
		server domain.Server
		err    error
	)
	data := DataResult{true, "ok", nil, nil}

	if readJson(r, "data", &server) {
		err = this.ServerRepo.Remove(server.Id)

		if err != nil {
			log.Printf("%v", err)
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (this DataHandler) loadServers(w http.ResponseWriter, r *http.Request) {
	var results []domain.Server
	this.ServerRepo.All(&results)

	data := DataResult{true, "ok", results, nil}
	json.NewEncoder(w).Encode(data)
}

func (this DataHandler) loadSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.XtvSettings
	this.SettingsRepo.FindFirst(&settings)

	data := DataResult{true, "ok", nil, &settings}
	json.NewEncoder(w).Encode(data)
}

func (this DataHandler) saveSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.XtvSettings
	var err error
	data := DataResult{true, "ok", nil, nil}

	if readJson(r, "data", &settings) {
		_, err = this.SettingsRepo.Save(settings.Id, &settings)

		if err != nil {
			log.Printf("%v", err)
			data.Success = false
			data.Status = "error"
		}
		this.DccService.UpdateSettings(&settings)
	}

	json.NewEncoder(w).Encode(data)
}
