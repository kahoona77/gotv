package handler

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/irc"
	"net/http"
)

type IrcHandler struct {
	Client *irc.IrcClient
}

func NewIrcHandler(client *irc.IrcClient) *IrcHandler {
	h := new(IrcHandler)
	h.Client = client
	return h
}

func (this IrcHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	switch {
	case url == "/irc/toggleConnection":
		this.toggleConnection(w, r)
	case url == "/irc/getServerStatus":
		this.getServerStatus(w, r)
	case url == "/irc/getServerConsole":
		this.getServerConsole(w, r)
	}

	return
}

func (this *IrcHandler) toggleConnection(w http.ResponseWriter, r *http.Request) {
	var server domain.Server
	data := map[string]interface{}{
		"success": true,
		"status":  "ok",
	}
	if readJson(r, "data", &server) {
		this.Client.ToggleConnection(&server)
		data["result"] = server
	}

	json.NewEncoder(w).Encode(data)
}

func (this *IrcHandler) getServerStatus(w http.ResponseWriter, r *http.Request) {
	var server domain.Server
	data := map[string]interface{}{
		"success": true,
		"status":  "undefined",
	}
	if readJson(r, "data", &server) {
		this.Client.GetServerStatus(&server)
		data["status"] = server.Status
	}

	json.NewEncoder(w).Encode(data)
}

func (this *IrcHandler) getServerConsole(w http.ResponseWriter, r *http.Request) {
	var server domain.Server
	data := map[string]interface{}{
		"success": true,
		"console": "",
	}
	if readJson(r, "data", &server) {
		data["console"] = this.Client.GetServerConsole(&server)
	}

	json.NewEncoder(w).Encode(data)
}
