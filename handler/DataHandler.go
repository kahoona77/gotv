package handler

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"log"
	"net/http"
)

type DataHandler struct {
	ServerRepo *domain.GoTvRepository
}

func NewDataHandler(serverRepo *domain.GoTvRepository) *DataHandler {
	h := new(DataHandler)
	h.ServerRepo = serverRepo
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
	case url == "/data/addChannel":
		this.addChannel(w, r)
	}

	return
}

func (this DataHandler) saveServer(w http.ResponseWriter, r *http.Request) {
	var (
		server domain.Server
		err    error
	)
	data := map[string]interface{}{
		"success": true,
		"status":  "ok",
	}
	if readJson(r, "data", &server) {

		_, err = this.ServerRepo.Save(server.Id, &server)

		if err != nil {
			log.Printf("%v", err)
		} else {
			data["success"] = true
			data["status"] = "ok"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (this DataHandler) deleteServer(w http.ResponseWriter, r *http.Request) {
	var (
		server domain.Server
		err    error
	)
	data := map[string]interface{}{
		"success": true,
		"status":  "ok",
	}
	if readJson(r, "data", &server) {
		err = this.ServerRepo.Remove(server.Id)

		if err != nil {
			log.Printf("%v", err)
		} else {
			data["success"] = true
			data["status"] = "ok"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (this DataHandler) addChannel(w http.ResponseWriter, r *http.Request) {
	data := map[string]*json.RawMessage{}
	result := map[string]interface{}{
		"success": true,
		"status":  "ok",
	}

	if readJson(r, "data", &data) {
		server := new(domain.Server)
		var (
			serverId string
		)
		json.Unmarshal(*data["serverId"], &serverId)
		this.ServerRepo.FindById(serverId, server)

		channel := new(domain.Channel)
		json.Unmarshal(*data["channel"], &channel)
		server.Channels = append(server.Channels, *channel)

		var (
			err error
		)
		_, err = this.ServerRepo.Save(server.Id, server)

		if err != nil {
			log.Printf("%v", err)
		} else {
			result["success"] = true
			result["status"] = "ok"
		}
	}

	json.NewEncoder(w).Encode(result)
}

func (this DataHandler) deleteChannel(w http.ResponseWriter, r *http.Request) {
	data := map[string]*json.RawMessage{}
	result := map[string]interface{}{
		"success": true,
		"status":  "ok",
	}

	if readJson(r, "data", &data) {
		server := new(domain.Server)
		var (
			serverId  string
			channelId string
		)
		json.Unmarshal(*data["serverId"], &serverId)
		this.ServerRepo.FindById(serverId, server)

		channel := new(domain.Channel)
		json.Unmarshal(*data["channelId"], &channelId)
		server.Channels = append(server.Channels, *channel)

		var (
			err error
		)
		_, err = this.ServerRepo.Save(server.Id, server)

		if err != nil {
			log.Printf("%v", err)
		} else {
			result["success"] = true
			result["status"] = "ok"
		}
	}

	json.NewEncoder(w).Encode(result)
}

func removeChannel(channels *[]domain.Channel, channel domain.Channel) {
	var i int
	for index, c := range *channels {
		if c == channel {
			i = index
		}
	}
	*channels = append(*channels[:i], *channels[i+1:]...)
}

func (this DataHandler) loadServers(w http.ResponseWriter, r *http.Request) {
	var results []domain.Server
	this.ServerRepo.All(&results)

	data := map[string]interface{}{
		"success": true,
		"status":  "ok",
		"results": results,
	}
	json.NewEncoder(w).Encode(data)
}
