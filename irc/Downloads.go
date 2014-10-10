package irc

import (
  "github.com/kahoona77/gotv/domain"
)

type Download struct {
  Id            string `json:"id"`
  Status        string `json:"status"`
  File          string `json:"file"`
  PacketId      string `json:"packetId"`
  Server        string `json:"server"`
  Bot           string `json:"bot"`
  BytesReceived int64  `json:"bytesReceived"`
  Size          int64  `json:"size"`
  Speed         float32  `json:"speed"`
  Remaining     int64  `json:"remaining"`
}

func DownloadFromPacket(packet *domain.Packet) *Download {
  d := Download {Id: packet.Name, Status: "WAITING", File: packet.Name, PacketId: packet.PacketId, Bot: packet.Bot, Server: packet.Server}
  return &d
}

func (dcc *DccService) DownloadPacket (packet *domain.Packet) {
  bot:= dcc.client.GetBot (packet.Server)
  bot.DownloadPacket (packet)
  download := DownloadFromPacket (packet)
  dcc.downloads[download.Id] = download
}

func (dcc *DccService) ListDownloads () []*Download {
  v := make([]*Download, 0, len(dcc.downloads))

  for  _, value := range dcc.downloads {
     v = append(v, value)
  }
  return v
}

func (dcc *DccService) updateDownload (downloadId string, totalBytes int64) {
  download := dcc.downloads[downloadId]
  download.Status = "RUNNING"
  download.BytesReceived = totalBytes

  //calc speed
  // sizeDelta := (totalBytes - download.BytesReceived) / 1024
  // timeDelta = (newTime - oldTime) / 1000
  //   return sizeDelta / timeDelta
}
