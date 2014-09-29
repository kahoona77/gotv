package domain

import (
 "time"
)

type MongoDomain interface {
    SetId(id string)
}

// Server
type Server struct {
	Id       string     `json:"id" bson:"_id"`
	Name     string     `json:"name" bson:"name"`
	Port     int        `json:"port" bson:"port"`
	Status   string     `json:"status" bson:"status"`
	Channels []Channel  `json:"channels" bson:"channels"`
}

func (this *Server) SetId(id string) {
	this.Id = id
}


type Channel struct {
	Name   string        `json:"name" bson:"name"`
}

//Packet
type Packet struct {
  Id        string     `json:"id" bson:"_id"`
  PacketId  string     `json:"packetId" bson:"packetId"`
  Size      string     `json:"size" bson:"size"`
  Name      string     `json:"name" bson:"name"`
  Bot       string     `json:"bot" bson:"bot"`
  Channel   string     `json:"channel" bson:"channel"`
  Server    string     `json:"server" bson:"server"`
  Date      time.Time  `json:"date" bson:"date"`
}

func (this *Packet) SetId(id string) {
  this.Id = id
}


// Settings
type XtvSettings struct {
  Id                   string     `json:"id" bson:"_id"`
  Nick                 string     `json:"nick" bson:"nick"`
  TempDir              string     `json:"tempDir" bson:"tempDir"`
  DownloadDir          string     `json:"downloadDir" bson:"downloadDir"`
  PostDownloadTrigger  string     `json:"postDownloadTrigger" bson:"postDownloadTrigger"`
  LogFile              string     `json:"logFile" bson:"logFile"`
  MaxDownStream        int        `json:"maxDownStream" bson:"maxDownStream"`
}

func (this *XtvSettings) SetId(id string) {
  this.Id = id
}
