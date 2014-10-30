package service

import (
  "log"
  "github.com/kahoona77/gotv/domain"
)

type Context struct {
  //DB
  MongoService *MongoService

  //Repositores
  ServerRepo    *GoTvRepository
  SettingsRepo  *GoTvRepository
  PacketsRepo   *GoTvRepository
  ShowsRepo     *GoTvRepository

  //Services
  DataService  *DataService
  DccService   *DccService
  IrcClient    *IrcClient
  ShowService  *ShowService
}

func (ctx *Context) GetSettings () *domain.XtvSettings{
  return ctx.DataService.GetSettings ()
}

func (ctx *Context) Close (){
  ctx.MongoService.Close()
}

func (ctx *Context) handleError (err error){
  if (err != nil) {
    log.Printf("ERROR: %v", err)
  }
}

func CreateContext (logFile string) *Context {
  c := new(Context)
  c.MongoService = CreateMongoService ()

  //Repositories
  c.ServerRepo   = NewRepository(c.MongoService.Session, "servers")
  c.SettingsRepo = NewRepository(c.MongoService.Session, "settings")
  c.PacketsRepo  = NewRepository(c.MongoService.Session, "packets")
  c.ShowsRepo    = NewRepository(c.MongoService.Session, "shows")

  //Service
  c.DataService = &DataService {c}
  c.DccService  = NewDccService (c)
  c.IrcClient   = NewIrcClient (c)
  c.ShowService = NewShowService (c)

  //set logfile in settings
  settings := c.GetSettings ()
  settings.LogFile = logFile
  c.DataService.SaveSettings (settings)

  return c
}
