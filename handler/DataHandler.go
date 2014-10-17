package handler

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/irc"
	"io/ioutil"
	"log"
	"net/http"
)

type DataResult struct {
	Success  bool                `json:"success"`
	Status   string              `json:"status"`
	Servers  []domain.Server     `json:"servers,omitempty"`
	Settings *domain.XtvSettings `json:"settings,omitempty"`
	LogFile  string              `json:"logFile,omitempty"`
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

func (dh DataHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	switch {
	case url == "/data/loadServers":
		dh.loadServers(w, r)
	case url == "/data/saveServer":
		dh.saveServer(w, r)
	case url == "/data/deleteServer":
		dh.deleteServer(w, r)
	case url == "/data/loadSettings":
		dh.loadSettings(w, r)
	case url == "/data/saveSettings":
		dh.saveSettings(w, r)
	case url == "/data/loadLogFile":
		dh.loadLogFile(w, r)
	case url == "/data/clearLogFile":
		dh.clearLogFile(w, r)
	}

	return
}

func (dh DataHandler) saveServer(w http.ResponseWriter, r *http.Request) {
	var (
		server domain.Server
		err    error
	)
	data := DataResult{true, "ok", nil, nil, ""}
	if readJson(r, "data", &server) {

		_, err = dh.ServerRepo.Save(server.Id, &server)

		if err != nil {
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (dh DataHandler) deleteServer(w http.ResponseWriter, r *http.Request) {
	var (
		server domain.Server
		err    error
	)
	data := DataResult{true, "ok", nil, nil, ""}

	if readJson(r, "data", &server) {
		err = dh.ServerRepo.Remove(server.Id)

		if err != nil {
			log.Printf("%v", err)
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (dh DataHandler) loadServers(w http.ResponseWriter, r *http.Request) {
	var results []domain.Server
	dh.ServerRepo.All(&results)

	data := DataResult{true, "ok", results, nil, ""}
	json.NewEncoder(w).Encode(data)
}

func (dh DataHandler) loadSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.XtvSettings
	dh.SettingsRepo.FindFirst(&settings)

	data := DataResult{true, "ok", nil, &settings, ""}
	json.NewEncoder(w).Encode(data)
}

func (dh DataHandler) saveSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.XtvSettings
	var err error
	data := DataResult{true, "ok", nil, nil, ""}

	if readJson(r, "data", &settings) {
		_, err = dh.SettingsRepo.Save(settings.Id, &settings)

		if err != nil {
			log.Printf("%v", err)
			data.Success = false
			data.Status = "error"
		}
		dh.DccService.UpdateSettings(&settings)
	}

	json.NewEncoder(w).Encode(data)
}

func (dh DataHandler) loadLogFile(w http.ResponseWriter, r *http.Request) {
	settings := dh.DccService.GetSettings()
	buf, _ := ioutil.ReadFile(settings.LogFile)

	data := DataResult{true, "ok", nil, nil, string(buf)}
	json.NewEncoder(w).Encode(data)
}

func (dh DataHandler) clearLogFile(w http.ResponseWriter, r *http.Request) {
	settings := dh.DccService.GetSettings()
	ioutil.WriteFile(settings.LogFile, []byte(""), 0644)

	data := DataResult{true, "ok", nil, nil, ""}
	json.NewEncoder(w).Encode(data)
}
