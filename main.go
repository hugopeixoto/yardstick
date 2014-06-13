package main

import (
  "time"
  "log"
  "os"
  "flag"
  "encoding/json"
  "strings"
  "io/ioutil"
  "github.com/peterbourgon/g2s"
)

type Config struct {
  Listen  string
  Statsd  string
  Peers   []Server
}

func Nop(Server, time.Duration, error) {
}

type Statsd struct {
  Statter g2s.Statter
}

func (s *Statsd) Report(endpoint Server, rtt time.Duration, err error) {
  hostname, _ := os.Hostname()

  if err != nil {
    log.Printf("%v error", hostname)
    s.Statter.Timing(1.0, hostname + ".ping." + strings.Split(endpoint.Name, ".")[0] + ".took", 3 * time.Second)
  } else {
    s.Statter.Timing(1.0, hostname + ".ping." + strings.Split(endpoint.Name, ".")[0] + ".took", rtt)
  }
}

func main () {
  var conf Config

  fn := Nop

  confpath := flag.String("conf", "servers.json", "")

  flag.Parse()

  file, _ := ioutil.ReadFile(*confpath)

  json.Unmarshal(file, &conf)

  if conf.Statsd != "" {
    e, _ := g2s.Dial("udp", conf.Statsd)
    fn = (&Statsd{e}).Report
  }

  NewYardstick(conf.Listen, conf.Peers, fn).Run()

  <-make(chan struct{})
}
