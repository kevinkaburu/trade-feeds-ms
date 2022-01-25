package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) TradeChat(w http.ResponseWriter, r *http.Request) {
	//LiveChat
	var response models.HttpResponse

	c, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Unable to upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var liveChat models.LiveChat
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("unable to read socket:", err)
			break
		}

		err = json.Unmarshal(message, &liveChat)
		if err != nil {
			log.Println("Unable to parse  OfferQuery Json  because ", err)
			response.Message = "Non json Message sent"
			response.Status = "400"
			c.WriteJSON(response)
		}
		//validate input

		if liveChat.TradeID <= 0 {
			log.Println("Trade ID not set", err)
			response.Message = "TradeID not set"
			response.Status = "400"
			c.WriteJSON(response)
		}
		//validate token
		//fetch user details from redis/Auth data

		//select Trade Data

		response.Message = "success"
		response.Status = "200"
		response.Data = liveChat
		c.WriteJSON(response)

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
