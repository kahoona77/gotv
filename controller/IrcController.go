package controller

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/service"
	"net/http"
)

type IrcController struct {
	*service.Context
}


func (ic IrcController) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	switch {
	case url == "/irc/toggleConnection":
		ic.toggleConnection(w, r)
	case url == "/irc/getServerStatus":
		ic.getServerStatus(w, r)
	case url == "/irc/getServerConsole":
		ic.getServerConsole(w, r)
	}

	return
}

func (ic *IrcController) toggleConnection(w http.ResponseWriter, r *http.Request) {
	var server domain.Server
	data := map[string]interface{}{
		"success": true,
		"status":  "ok",
	}
	if readJson(r, "data", &server) {
		ic.IrcClient.ToggleConnection(&server)
		data["result"] = server
	}

	json.NewEncoder(w).Encode(data)
}

func (ic *IrcController) getServerStatus(w http.ResponseWriter, r *http.Request) {
	var server domain.Server
	data := map[string]interface{}{
		"success": true,
		"status":  "undefined",
	}
	if readJson(r, "data", &server) {
		ic.IrcClient.GetServerStatus(&server)
		data["status"] = server.Status
	}

	json.NewEncoder(w).Encode(data)
}

func (ic *IrcController) getServerConsole(w http.ResponseWriter, r *http.Request) {
	var server domain.Server
	data := map[string]interface{}{
		"success": true,
		"console": "",
	}
	if readJson(r, "data", &server) {
		data["console"] = ic.IrcClient.GetServerConsole(&server)
	}

	json.NewEncoder(w).Encode(data)
}
