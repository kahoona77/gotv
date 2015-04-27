package domain

import (
	"time"
)

type MongoDomain interface {
	SetId(id string)
}

// Server
type Server struct {
	Id       string    `json:"id" bson:"_id"`
	Name     string    `json:"name" bson:"name"`
	Port     int       `json:"port" bson:"port"`
	Status   string    `json:"status" bson:"status"`
	Channels []Channel `json:"channels" bson:"channels"`
}

func (this *Server) SetId(id string) {
	this.Id = id
}

type Channel struct {
	Name string `json:"name" bson:"name"`
}

//Packet
type Packet struct {
	Id       string    `json:"id" bson:"_id"`
	PacketId string    `json:"packetId" bson:"packetId"`
	Size     string    `json:"size" bson:"size"`
	Name     string    `json:"name" bson:"name"`
	Bot      string    `json:"bot" bson:"bot"`
	Channel  string    `json:"channel" bson:"channel"`
	Server   string    `json:"server" bson:"server"`
	Date     time.Time `json:"date" bson:"date"`
}

func NewPacket(packetId string, size string, name string, bot string, channel string, server string, date time.Time) *Packet {
	p := new(Packet)
	p.Id = channel + ":" + bot + ":" + packetId
	p.PacketId = packetId
	p.Size = size
	p.Name = name
	p.Bot = bot
	p.Channel = channel
	p.Server = server
	p.Date = date
	return p
}

func (this *Packet) SetId(id string) {
	this.Id = id
}

// Settings
type XtvSettings struct {
	Id            string `json:"id" bson:"_id"`
	Nick          string `json:"nick" bson:"nick"`
	TempDir       string `json:"tempDir" bson:"tempDir"`
	DownloadDir   string `json:"downloadDir" bson:"downloadDir"`
	ShowsFolder   string `json:"showsFolder" bson:"showsFolder"`
	MoviesFolder  string `json:"moviesFolder" bson:"moviesFolder"`
	KodiAddress   string `json:"kodiAddress" bson:"kodiAddress"`
	TraktToken    string `json:"traktToken" bson:"traktToken"`
	LogFile       string `json:"logFile" bson:"logFile"`
	MaxDownStream int    `json:"maxDownStream" bson:"maxDownStream"`
}

func (this *XtvSettings) SetId(id string) {
	this.Id = id
}

//Show
type Show struct {
	Id         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	Banner     string `json:"banner" bson:"banner"`
	FirstAired string `json:"firstAired" bson:"firstAired"`
	Overview   string `json:"overview" bson:"overview"`
	SearchName string `json:"searchName" bson:"searchName"`
	Folder     string `json:"folder" bson:"folder"`
}

func (this *Show) SetId(id string) {
	this.Id = id
}

//Episode
type Episode struct {
	ID            string `json:"id" bson:"_id"`
	Name          string `json:"name" bson:"name"`
	FirstAired    string `json:"firstAired" bson:"firstAired"`
	Overview      string `json:"overview" bson:"overview"`
	Filename      string `json:"filename" bson:"filename"`
	EpisodeNumber uint64 `json:"episodeNumber" bson:"episodeNumber"`
	SeasonNumber  uint64 `json:"seasonNumber" bson:"seasonNumber"`
}

func (this *Episode) SetId(id string) {
	this.ID = id
}
