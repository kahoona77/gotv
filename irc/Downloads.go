package irc

import (
  "log"
  "time"
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
  LastUpdate    time.Time  `json:"-"`
}

func DownloadFromPacket(packet *domain.Packet) *Download {
  d := Download {Id: packet.Name, Status: "WAITING", File: packet.Name, PacketId: packet.PacketId, Bot: packet.Bot, Server: packet.Server}
  d.LastUpdate = time.Now()
  return &d
}

func (dcc *DccService) DownloadPacket (packet *domain.Packet) {
  bot:= dcc.client.GetBot (packet.Server)
  download := DownloadFromPacket (packet)
  dcc.downloads[download.Id] = download
  bot.StartDownload (download)
}

func (dcc *DccService) ListDownloads () []*Download {
  v := make([]*Download, 0, len(dcc.downloads))

  for  _, value := range dcc.downloads {
     v = append(v, value)
  }
  return v
}

func (dcc *DccService) StopDownload (download *Download) {
  bot:= dcc.client.GetBot (download.Server)
  bot.StopDownload (download)
}

func (dcc *DccService) CancelDownload (parsedDownload *Download) {
  download := dcc.downloads[parsedDownload.Id]
  if (download != nil) {
    if (download.Status == "RUNNING") {
      dcc.StopDownload (download)
    }
    delete (dcc.downloads, download.Id)
  }
}

func (dcc *DccService) ResumeDownload (parsedDownload *Download) {
  download := dcc.downloads[parsedDownload.Id]
  if (download != nil) {
    if (download.Status != "RUNNING") {
      dcc.StopDownload (download)
    }
  }
}

func (dcc *DccService) updateDownloads () {
  for {
    update := <- dcc.updateChan
    download := dcc.downloads[update.File]
    if (download != nil) {
      //calc speed
      now:= time.Now()
      sizeDelta := (update.TotalBytes - download.BytesReceived) / 1024
      timeDelta := (now.UnixNano() - download.LastUpdate.UnixNano())
      download.Speed = (float32(sizeDelta) / float32(timeDelta)) * 1000 * 1000 * 1000

      //update download
      download.LastUpdate = now
      download.Status = "RUNNING"
      download.BytesReceived = update.TotalBytes
      download.Size = update.Size
    } else {
      log.Printf("download not found: %v in %v", update.File, dcc.downloads)
    }
  }
}


func (dcc *DccService) completeDownload (file string) {
  download := dcc.downloads[file]
  if (download != nil) {
    log.Printf("Download completed '%v'", download.File)
    download.Status = "COMPLETE"

    //TODO move file to destination

  } else {
    log.Printf("download not found: %v in %v",file, dcc.downloads)
  }
}
