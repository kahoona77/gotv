package service

import (
	"github.com/kahoona77/gotv/domain"
	"strings"
)

type IrcClient struct {
	*Context
	bots         map[string]*IrcBot
}

func NewIrcClient(ctx *Context) *IrcClient {
	client := new(IrcClient)
	client.Context = ctx
	client.bots    = make(map[string]*IrcBot)
	return client
}

func (ic *IrcClient) ToggleConnection (server *domain.Server) {
	bot := ic.getAndUpdateBot (server)
	if (bot.IsConnected()){
		bot.Disconnect ()
		server.Status = "Not Connected"
	} else {
    bot.Connect()
    server.Status = "Connected"
  }
  ic.ServerRepo.Save(server.Id, server)
}


func (ic *IrcClient) GetServerStatus(server *domain.Server) {
  bot := ic.getAndUpdateBot (server)
  if (bot.IsConnected()) {
    server.Status = "Connected"
  } else {
    server.Status = "Not Connected"
  }
	ic.ServerRepo.Save(server.Id, server)
}

func (ic *IrcClient) GetServerConsole(server *domain.Server) string{
  bot := ic.getAndUpdateBot (server)
  return strings.Join (bot.ConsoleLog, "\n")
}

func (ic *IrcClient) getAndUpdateBot (server *domain.Server) *IrcBot{
	bot := ic.bots[server.Name]
	if (bot == nil) {
		// create new bot
		bot = NewIrcBot (ic.Context, server)
		ic.bots[server.Name] = bot
	} else {
		//update bot
		bot.Server   = server
	}
	return bot
}

func (ic *IrcClient) GetBot (serverName string) *IrcBot{
	return ic.bots[serverName]
}
