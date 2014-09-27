package domain

import (

)

type MongoDomain interface {
    SetId(id string)
}

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
