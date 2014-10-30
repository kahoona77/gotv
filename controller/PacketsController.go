package controller

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/service"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strings"
	"time"
)

type PacketsResult struct {
	Success bool            `json:"success"`
	Status  string          `json:"status"`
	Packets []domain.Packet `json:"packets,omitempty"`
	Count   int             `json:"count,omitempty"`
}

type PacketsController struct {
	*service.Context
}

func (pc PacketsController) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	switch {
	case strings.HasPrefix(url, "/packets/findPackets"):
		pc.findPackets(w, r)
	case url == "/packets/countPackets":
		pc.countPackets(w, r)
	}

	return
}

func (pc PacketsController) findPackets(w http.ResponseWriter, r *http.Request) {
	var results []domain.Packet

	params := r.URL.Query()
	queryRegex := createRegexQuery(params["query"][0])

	query := bson.M{"name": bson.M{"$regex": queryRegex, "$options": "i"}}

	pc.PacketsRepo.FindWithQuery(&query, &results)

	data := PacketsResult{true, "ok", results, 0}
	json.NewEncoder(w).Encode(data)
}

func createRegexQuery(query string) string {
	parts := strings.Split(query, " ")
	return strings.Join(parts, ".*")
}

func (pc PacketsController) countPackets(w http.ResponseWriter, r *http.Request) {
	//clean old packets
	//[date: ['$lt': yesterday.format("yyyy-MM-dd'T'HH:mm:ss.SSSZ")]]
	minusOneDay, _ := time.ParseDuration("-24h")
	yesterday := time.Now().Add(minusOneDay)
	removeQuery := bson.M{"date": bson.M{"$lt": yesterday}}

	info, err := pc.PacketsRepo.RemoveAll(&removeQuery)
	if err != nil {
		log.Printf("error while deleting old packets: %v", err)
	} else {
		log.Printf("removed %v old packets", info.Removed)
	}

	count, _ := pc.PacketsRepo.CountAll()

	data := PacketsResult{true, "ok", nil, count}
	json.NewEncoder(w).Encode(data)
}
