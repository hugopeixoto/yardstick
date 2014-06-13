package main

import (
  "net"
  "fmt"
  "time"
  "bytes"
  "crypto/rand"
)

type Server struct {
  Name    string
  Address string
}

type Yardstick struct {
  Server    Server
  Endpoints []Server
}

func NewYardstick(name string, servers []Server) (*Yardstick) {
  nf := Yardstick{}

  nf.Endpoints = make([]Server, 0)

  for _, server := range servers {
    if server.Name == name {
      nf.Server = server
    } else {
      nf.Endpoints = append(nf.Endpoints, server)
    }
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

  listenAddress, _ := net.ResolveUDPAddr("udp", nw.Server.Address)
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

  connection, _ := net.Dial("udp", endpoint.Address)

  before := time.Now()
  connection.Write(ping)

  connection.SetReadDeadline(time.Now().Add(5 * time.Second))

  nbytes, err := connection.Read(pong)
  after := time.Now()

  if err != nil || bytes.Compare(pong[0:nbytes], ping) != 0 {
    fmt.Printf("<%v> %v Error\n", nw.Server.Name, endpoint.Name)
  } else {
    fmt.Printf("<%v> %v RTT: %v\n", nw.Server.Name, endpoint.Name, after.Sub(before))
  }
}
