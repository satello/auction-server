package main

import (
  "log"
)

type Tower struct {
  rooms map[string]*DraftHub
}

func newTower() *Tower {
  return &Tower{
    rooms: make(map[string]*DraftHub),
  }
}

type Room struct {
  RoomId string `json:"roomId"`

  Rules *Rules `json:"rules"`

  Bidders []*Bidder `json:"bidders"`

  Players []*Player `json:"players"`
}

func newRoom(t *Tower, rules *Rules, bidders []*Bidder, players []*Player) string {
  log.Println("starting new room")
  log.Printf("number of bidders %d", len(bidders))
  log.Println(bidders)
  roomId := createUuid()
  newDraftRoom := newDraft(rules, bidders, players)

  // start new hub
  go newDraftRoom.run()

  // TODO watch out or memory leaks with this. Do go routines shut down when the parent does?
  log.Println("new draft room created")
  log.Println(roomId)
  log.Printf("number of bidders %d", len(newDraftRoom.bidders))
  t.rooms[roomId] = newDraftRoom

  return roomId
}
