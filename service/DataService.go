package service

import (
  "github.com/kahoona77/gotv/domain"
)

type DataService struct {
  *Context
}

func (ds *DataService)  FindAllServers () ([]domain.Server, error) {
  var results []domain.Server
  err := ds.ServerRepo.All(&results)
  ds.handleError(err)
  return results, err
}

func (ds *DataService)  DeleteServer (server *domain.Server) error {
  err := ds.ServerRepo.Remove(server.Id)
  ds.handleError(err)
  return err
}

func (ds *DataService)  SaveServer (server *domain.Server) error {
  _, err := ds.ServerRepo.Save(server.Id, server)
  ds.handleError(err)
  return err
}

func (ds *DataService) GetSettings () *domain.XtvSettings {
  var settings domain.XtvSettings
  ds.SettingsRepo.FindFirst(&settings)
  return &settings
}

func (ds *DataService)  SaveSettings (settings *domain.XtvSettings) error {
  _, err := ds.SettingsRepo.Save(settings.Id, settings)
  ds.handleError(err)
  return err
}
