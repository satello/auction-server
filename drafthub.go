package main

import (
  "fmt"
  "encoding/json"
  "log"
  "github.com/mitchellh/mapstructure"
)

// draft hub maintains the set of active clients and broadcasts messages to the
// clients.
type DraftHub struct {
	// Registered clients.
	clients map[*Subscriber]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Subscriber

	// Unregister requests from clients.
	unregister chan *Subscriber

  // accept message from client
  acceptMessage chan *Message

  // players eligable for draft
  players map[string]*Player

  // bidders in the draft
  bidders map[string]*Bidder

  // draft rules
  rules *Rules

  // flag to set when you want to close draft room
  isActive bool
}

func newDraft(rules *Rules, bidders []*Bidder, players []*Player) *DraftHub {
  bidder_map := make(map[string]*Bidder)
  for _, v := range bidders {
    log.Printf("%s", v.BidderId)
    bidder_map[v.BidderId] = v
  }

	return &DraftHub{
		broadcast:      make(chan []byte),
		register:       make(chan *Subscriber),
		unregister:     make(chan *Subscriber),
    acceptMessage:  make(chan *Message),
		clients:        make(map[*Subscriber]bool),
    players:        make(map[string]*Player),
    bidders:        bidder_map,
    rules:          rules,
    isActive:       true,
	}
}

func (h *DraftHub) run() {
	for {
		select {

		case client := <-h.register:
      log.Println("CONNECTING CLIENT")
			h.clients[client] = true

		case client := <-h.unregister:
      log.Println("DISCONNECTING CLIENT")
			if _, ok := h.clients[client]; ok {
        delete(h.bidders, client.bidderId)
				delete(h.clients, client)
				close(client.send)
			}

		case messageJson := <-h.acceptMessage:
      switch t := messageJson.MessageType; t {

    	case "newBidder":
        var body NewBidderBody
        mapstructure.Decode(messageJson.Body, &body)
        name := body.Name
        cap := body.Cap
        spots := body.Spots
        createBidder(name, cap, spots, messageJson.Subscriber, h)

      case "authorizeBidder":
        var body TokenBody
        mapstructure.Decode(messageJson.Body, &body)

        token := body.Token
        authorizeBidder(token, messageJson.Subscriber, h)

      case "deauthorizeBidder":
        var body TokenBody
        mapstructure.Decode(messageJson.Body, &body)

        token := body.Token
        deactivateBidder(token, messageJson.Subscriber, h)

      case "getBidders":
        getBidders(messageJson.Subscriber, h)

    	case "chatMessage":
        log.Printf("CHAT MESSAGE");
        body := messageJson.Body

        response := Response{"CHAT_MESSAGE", body}
        response_json, err := json.Marshal(response)
        if err != nil {
    			log.Printf("error: %v", err)
    			break
        }
        broadcastMessage(h, response_json)

    	default:
    		// freebsd, openbsd,
    		// plan9, windows...
    		fmt.Printf("%s.", t)
      }
		}
	}
}
