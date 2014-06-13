package main

import (
  "net"
  "log"
  "time"
  "bytes"
  "crypto/rand"
  "errors"
)

type Server struct {
  Name    string
  Address string
}

type Yardstick struct {
  Address   string
  Endpoints []Server
}

func NewYardstick(listen string, servers []Server) (*Yardstick) {
  nf := Yardstick{}

  nf.Address   = listen
  nf.Endpoints = make([]Server, 0)

  for _, server := range servers {
    if server.Address == "" {
      server.Address = server.Name + ":4367"
    }

    nf.Endpoints = append(nf.Endpoints, server)
  }

  return &nf
}

func (nw *Yardstick) Run() {
  go nw.Listen()

  for {
    for _, server := range nw.Endpoints {
      go nw.Ping(server)
    }

    time.Sleep(3 * time.Second)
  }
}

func (nw *Yardstick) Listen() {
  message := make([]byte, 256)

  listenAddress, _ := net.ResolveUDPAddr("udp", nw.Address)
  sock, _          := net.ListenUDP("udp", listenAddress)

  for {
    bytes, peerAddress, _ := sock.ReadFromUDP(message)

    sock.WriteToUDP(message[0:bytes], peerAddress)
  }
}

func (nw *Yardstick) Ping(endpoint Server) {
  ping := make([]byte, 256)
  pong := make([]byte, 256)

  rand.Read(ping)

  connection, err := net.Dial("udp", endpoint.Address)
  if err != nil {
    nw.Report(endpoint, 0, err)
    return
  }

  before := time.Now()
  _, err = connection.Write(ping)

  if err != nil {
    nw.Report(endpoint, 0, err)
    return
  }

  connection.SetReadDeadline(time.Now().Add(3 * time.Second))

  nbytes, err := connection.Read(pong)
  after := time.Now()

  if err != nil {
    nw.Report(endpoint, 0, err)
    return
  }

  if bytes.Compare(pong[0:nbytes], ping) != 0 {
    nw.Report(endpoint, 0, errors.New("Received invalid nonce"))
    return
  }


  nw.Report(endpoint, before.Sub(after), nil)
}

func (nw *Yardstick) Report(endpoint Server, rtt time.Duration, err error) {
  if err != nil {
    log.Printf("<%v> %v Error: %v\n", nw.Address, endpoint.Name, err)
  } else {
    log.Printf("<%v> %v RTT: %v\n", nw.Address, endpoint.Name, rtt)
  }
}
