package irc

import (
  "log"
  "strings"
  "regexp"
  irc "github.com/fluffle/goirc/client"
)

func Connect() {
    // create a config and fiddle with it first:
    cfg := irc.NewConfig("kahoona-go")
    cfg.Server = "irc.abjects.net:6667"
    c := irc.Client(cfg)

    // Add handlers to do things here!
    // e.g. join a channel on connect.
    c.HandleFunc("connected",
        func(conn *irc.Conn, line *irc.Line) {
          log.Print ("connected")
          conn.Join("#mg-chat")
          conn.Join("#moviegods")
        })

    c.HandleFunc("PRIVMSG", parsePacket)
    // And a signal on disconnect
    quit := make(chan bool)
    c.HandleFunc("disconnected",
        func(conn *irc.Conn, line *irc.Line) { quit <- true })

    // Tell client to connect.
    log.Print ("connecting")
    if err := c.Connect(); err != nil {
        log.Printf("Connection error: %v\n", err)
    }

    // ... or, use ConnectTo instead.
    if err := c.ConnectTo("irc.freenode.net"); err != nil {
        log.Printf("Connection error: %v\n", err)
    }

    // Wait for disconnect
    <-quit
}



func parsePacket (conn *irc.Conn, line *irc.Line) {
  r, _ := regexp.Compile(`(#[0-9]+).*\[\s*([0-9|\.]+[BbGgiKMs]+)\]\s+(.+).*`)

  // matches:= r.FindAllStringSubmatch(line.Text(), -1)
  result:= r.FindStringSubmatch(line.Text())
  log.Printf ("id: %v; size: %v, packet: %v\n", result[1], result[2], strings.TrimSpace(result[3]))

  // for k, v := range result {
  //     log.Printf("%d. %s\n", k, v)
  // }
}
