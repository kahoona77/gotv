package controller

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/service"
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

type DataController struct {
	*service.Context
}

func (dc DataController) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	switch {
	case url == "/data/loadServers":
		dc.loadServers(w, r)
	case url == "/data/saveServer":
		dc.saveServer(w, r)
	case url == "/data/deleteServer":
		dc.deleteServer(w, r)
	case url == "/data/loadSettings":
		dc.loadSettings(w, r)
	case url == "/data/saveSettings":
		dc.saveSettings(w, r)
	case url == "/data/loadLogFile":
		dc.loadLogFile(w, r)
	case url == "/data/clearLogFile":
		dc.clearLogFile(w, r)
	}

	return
}

func (dc DataController) saveServer(w http.ResponseWriter, r *http.Request) {
	var (
		server domain.Server
		err    error
	)
	data := DataResult{true, "ok", nil, nil, ""}
	if readJson(r, "data", &server) {

		err = dc.DataService.SaveServer(&server)

		if err != nil {
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (dc DataController) deleteServer(w http.ResponseWriter, r *http.Request) {
	var (
		server domain.Server
		err    error
	)
	data := DataResult{true, "ok", nil, nil, ""}

	if readJson(r, "data", &server) {
		err = dc.DataService.DeleteServer(&server)

		if err != nil {
			log.Printf("%v", err)
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (dc DataController) loadServers(w http.ResponseWriter, r *http.Request) {
	results, _ := dc.DataService.FindAllServers()

	data := DataResult{true, "ok", results, nil, ""}
	json.NewEncoder(w).Encode(data)
}

func (dc DataController) loadSettings(w http.ResponseWriter, r *http.Request) {
	settings := dc.GetSettings ()

	data := DataResult{true, "ok", nil, settings, ""}
	json.NewEncoder(w).Encode(data)
}

func (dc DataController) saveSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.XtvSettings
	var err error
	data := DataResult{true, "ok", nil, nil, ""}

	if readJson(r, "data", &settings) {
		_, err = dc.DataService.SettingsRepo.Save(settings.Id, &settings)

		if err != nil {
			log.Printf("%v", err)
			data.Success = false
			data.Status = "error"
		}
	}

	//update DownloadLimit
	dc.DccService.SetDownloadLimit(settings.MaxDownStream)

	json.NewEncoder(w).Encode(data)
}

func (dc DataController) loadLogFile(w http.ResponseWriter, r *http.Request) {
	settings := dc.GetSettings()
	buf, _ := ioutil.ReadFile(settings.LogFile)

	data := DataResult{true, "ok", nil, nil, string(buf)}
	json.NewEncoder(w).Encode(data)
}

func (dc DataController) clearLogFile(w http.ResponseWriter, r *http.Request) {
	settings := dc.GetSettings()
	ioutil.WriteFile(settings.LogFile, []byte(""), 0644)

	data := DataResult{true, "ok", nil, nil, ""}
	json.NewEncoder(w).Encode(data)
}
