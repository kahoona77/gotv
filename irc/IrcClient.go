package irc

import (
	"github.com/kahoona77/gotv/domain"
	"strings"
)

type IrcClient struct {
	PacketsRepo  *domain.GoTvRepository
	ServerRepo   *domain.GoTvRepository
	DccService   *DccService
	Settings     *domain.XtvSettings
	bots         map[string]*IrcBot
}

func NewClient(packetsRepo *domain.GoTvRepository, serverRepo *domain.GoTvRepository, settings *domain.XtvSettings) *IrcClient {
	client := new(IrcClient)
	client.PacketsRepo  = packetsRepo
	client.ServerRepo   = serverRepo
	client.Settings     = settings
	client.bots         = make(map[string]*IrcBot)
	return client
}

func (this *IrcClient) ToggleConnection (server *domain.Server) {
	bot := this.getAndUpdateBot (server)
	if (bot.IsConnected()){
		bot.Disconnect ()
		server.Status = "Not Connected"
	} else {
    bot.Connect()
    server.Status = "Connected"
  }
  this.ServerRepo.Save(server.Id, server)
}


func (this *IrcClient) GetServerStatus(server *domain.Server) {
  bot := this.getAndUpdateBot (server)
  if (bot.IsConnected()) {
    server.Status = "Connected"
  } else {
    server.Status = "Not Connected"
  }
	this.ServerRepo.Save(server.Id, server)
}

func (this *IrcClient) GetServerConsole(server *domain.Server) string{
  bot := this.getAndUpdateBot (server)
  return strings.Join (bot.ConsoleLog, "\n")
}

func (this *IrcClient) getAndUpdateBot (server *domain.Server) *IrcBot{
	bot := this.bots[server.Name]
	if (bot == nil) {
		// create new bot
		bot = NewIrcBot (this.PacketsRepo, this.DccService, this.Settings, server)
		this.bots[server.Name] = bot
	} else {
		//update bot
		bot.Server   = server
		bot.Settings = this.Settings
	}
	return bot
}

func (this *IrcClient) GetBot (serverName string) *IrcBot{
	return this.bots[serverName]
}
