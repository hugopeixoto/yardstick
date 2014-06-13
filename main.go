package main

import (
  "flag"
  "encoding/json"
  "io/ioutil"
)

func main () {
  var servers []Server

  conf := flag.String("conf", "servers.json", "")
  name := flag.String("name", "", "")

  flag.Parse()

  file, _ := ioutil.ReadFile(*conf)

  json.Unmarshal(file, &servers)

  NewYardstick(*name, servers).Run()

  <-make(chan struct{})
}
