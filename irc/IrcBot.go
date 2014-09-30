package irc

import (
  irc "github.com/fluffle/goirc/client"
  "github.com/kahoona77/gotv/domain"
  "log"
  "strconv"
  "regexp"
)

type IrcBot struct {
  PacketsRepo *domain.GoTvRepository
  Settings    *domain.XtvSettings
  Server      *domain.Server
  Conn        *irc.Conn
  ConsoleLog  []string
  Regex       *regexp.Regexp
  LogCount    int
}

func NewIrcBot(packetsRepo *domain.GoTvRepository, settings *domain.XtvSettings, server *domain.Server) *IrcBot {
  bot := new(IrcBot)
  bot.PacketsRepo = packetsRepo
  bot.Settings    = settings
  bot.Server      = server
  bot.ConsoleLog  = make ([]string, 0)
  bot.Regex, _    = regexp.Compile(`(#[0-9]+).*\[\s*([0-9|\.]+[BbGgiKMs]+)\]\s+(.+).*`)
  return bot
}

func (this *IrcBot) IsConnected() bool {
  if (this.Conn == nil) {
    return false
  }
  return this.Conn.Connected()
}

func (this *IrcBot) Connect() {
  // create a config and fiddle with it first:
  cfg := irc.NewConfig(this.Settings.Nick)
  cfg.Server = this.Server.Name + ":" + strconv.Itoa(this.Server.Port)
  this.Conn = irc.Client(cfg)

  // Join channels
  this.Conn.HandleFunc("connected",
    func(conn *irc.Conn, line *irc.Line) {
      this.logToConsole("connected to " + this.Server.Name + ":" + strconv.Itoa(this.Server.Port))

      for _, channel := range this.Server.Channels {
        this.logToConsole("joining channel " + channel.Name)
        conn.Join(channel.Name)
      }
    })

  // Parse Messages
  this.Conn.HandleFunc("PRIVMSG", this.parseMessage)


  // Tell client to connect.
  log.Print("connecting")
  if err := this.Conn.Connect(); err != nil {
    log.Printf("Connection error: %v\n", err)
    this.logToConsole("Connection error: " + err.Error())
  }
}

func (this *IrcBot) Disconnect() {
  // TODO
  //this.Conn.shutdown()
}

func (this *IrcBot) parseMessage (conn *irc.Conn, line *irc.Line) {
  packet := this.parsePacket (conn, line)
  if (packet == nil) {
    // log message
    this.logToConsole (line.Text())
  }
}

func (this *IrcBot) parsePacket(conn *irc.Conn, line *irc.Line) *domain.Packet {
  result := this.Regex.FindStringSubmatch(line.Text())
  packet := domain.NewPacket (result[1], result[2], result[3], line.Nick, line.Target(), this.Server.Name, line.Time)

  //save packet
  this.PacketsRepo.Save(packet.Id, packet)

  return packet
}

func (this *IrcBot) logToConsole (msg string) {
  if (this.LogCount > 500) {
    this.ConsoleLog = make ([]string, 0)
    this.LogCount = 0
  }
  this.ConsoleLog = append (this.ConsoleLog, msg)
  this.LogCount++
}
