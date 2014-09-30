package handler

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"log"
	"strings"
	"net/http"
	"labix.org/v2/mgo/bson"
)

type PacketsResult struct {
	Success  bool                `json:"success"`
	Status   string              `json:"status"`
	Packets  []domain.Packet     `json:"packets,omitempty"`
	Count    int                 `json:"count,omitempty"`
}

type PacketsHandler struct {
	PacketsRepo   *domain.GoTvRepository
}

func NewPacketsHandler(packetsRepo *domain.GoTvRepository) *PacketsHandler {
	h := new(PacketsHandler)
	h.PacketsRepo = packetsRepo
	return h
}

func (this PacketsHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.Print("URL: " + url)

	switch {
	case strings.HasPrefix(url, "/packets/findPackets"):
		this.findPackets(w, r)
	case url == "/packets/countPackets":
		this.countPackets(w, r)
	}

	return
}

func (this PacketsHandler) findPackets(w http.ResponseWriter, r *http.Request) {
	var results []domain.Packet

	params := r.URL.Query()
	queryRegex := createRegexQuery(params["query"][0])

	query := bson.M{"name": bson.M{"$regex": queryRegex, "$options": "i"}}


	this.PacketsRepo.FindWithQuery(&query, &results)

	data := PacketsResult {true,"ok", results, 0}
	json.NewEncoder(w).Encode(data)
}

func createRegexQuery (query string) string {
    parts := strings.Split (query, " ")
    return  strings.Join (parts, ".*")
}


func (this PacketsHandler) countPackets(w http.ResponseWriter, r *http.Request) {
	count, _ := this.PacketsRepo.CountAll()

	data := PacketsResult {true,"ok", nil, count}
	json.NewEncoder(w).Encode(data)
}
