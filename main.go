package main

import (
  "flag"
  "encoding/json"
  "io/ioutil"
)

type Config struct {
  Listen  string
  Peers   []Server
}

func main () {
  var conf Config

  confpath := flag.String("conf", "servers.json", "")

  flag.Parse()

  file, _ := ioutil.ReadFile(*confpath)

  json.Unmarshal(file, &conf)

  NewYardstick(conf.Listen, conf.Peers).Run()

  <-make(chan struct{})
}
