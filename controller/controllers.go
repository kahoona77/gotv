package controller

import (
  "encoding/json"
  "log"
  "net/http"
)


func readJson(req *http.Request, field string, v interface{}) bool {
  defer req.Body.Close()
  decoder := json.NewDecoder(req.Body)


  result := map[string]*json.RawMessage {}

  err := decoder.Decode(&result)
  if err != nil {
    log.Printf("ReadJson couldn't read request body %v", err)
    return false
  }

  json.Unmarshal(*result[field], &v)

  return true
}
