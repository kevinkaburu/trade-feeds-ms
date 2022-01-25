package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
		var (
			TradeID    uint64
			ProfileID  uint64
			ExternalID string
		)
		//fetch Offer from DB
		selectOfferQuery := fmt.Sprintf("select t.trade_id,t.profile_id,t.external_id from trade t where trade_id =%v;", liveChat.TradeID)
		err = s.DB.QueryRow(selectOfferQuery).Scan(&TradeID, &ProfileID, &ExternalID)
		if err != nil {
			log.Println("Unknown Trade  required")
			response.Message = "Trade  not found."
			response.Status = "400"
			c.WriteJSON(response)
		}
		//if message available send it
		if len(liveChat.Message) > 0 {
			if !s.SendMessage(ExternalID, liveChat.Message) {
				log.Println("Unable to send message to paxful seller.")
				response.Message = "Unable to send message to seller. Try again"
				response.Status = "400"
				c.WriteJSON(response)
			}
		}

		//fetch Messages
		pticker := time.NewTicker(10 * time.Second)
		pquit := make(chan struct{})
		func() {
			for {
				select {
				case <-pticker.C:
					s.refreshMessages(ExternalID, TradeID, c)

				case <-pquit:
					pticker.Stop()
					return
				}
			}
		}()

		// response.Message = "success"
		// response.Status = "200"
		// response.Data = messages
		// c.WriteJSON(response)

		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

func (s *Server) refreshMessages(ExternaltradeID string, tradeID uint64, c *websocket.Conn) (newTradeResponse models.StartTradeResponse) {

	var (
		ChatID     uint64
		Type       string
		ExternalID string
		Author     string
		Vuser      uint64
		Timestamp  uint64
	)
	//fetch Offer from DB
	selectOfferQuery := fmt.Sprintf("select trade_chat_id,type,external_id,author,villager_user,timestamp from trade_chat where trade_id =%v order by 1 desc limit 1;", tradeID)
	err := s.DB.QueryRow(selectOfferQuery).Scan(&ChatID, &Type, &ExternalID, &Author, &Vuser, &Timestamp)
	if err != nil {
		log.Println("Reading chats from DB:", err)
	}
	allchats := s.FetchMessages(ExternaltradeID)

	//loop through them inserting into db.
	lastMessage := false
	for i := 0; i < len(allchats.Data.Messages); i++ {
		if !lastMessage && allchats.Data.Messages[i].ID != ExternalID {
			continue
		} else if !lastMessage && allchats.Data.Messages[i].ID == ExternalID {
			lastMessage = true
		}
		var trade_chat_id uint64
		selectOfferQuery := fmt.Sprintf("select trade_chat_id from trade_chat where trade_id =%v and external_id = '%v' order by 1 desc limit 1;", tradeID, allchats.Data.Messages[i].ID)
		err := s.DB.QueryRow(selectOfferQuery).Scan(&trade_chat_id)
		if err != nil {
			log.Println("Reading chats from DB:", err)
		}
		if trade_chat_id > 0 {
			continue
		}

		var chats models.ChatMessages
		InsertChatQuery := "insert ignore into trade_chat set trade_id=?,type=?,external_id=?,author=?,text=?,villager_user=?,timestamp=?,created=now(),modified=now()"
		villagerUser := 0
		if allchats.Data.Messages[i].Author == os.Getenv("VILLAGER_PAXFUL_USERNAME") {
			villagerUser = 1
		}
		smsObject, err := s.DB.Exec(InsertChatQuery, tradeID, allchats.Data.Messages[i].Type, allchats.Data.Messages[i].ID, allchats.Data.Messages[i].Author, allchats.Data.Messages[i].Text, villagerUser, allchats.Data.Messages[i].Timestamp)
		if err != nil {
			log.Printf("PAXFUL insert ignore trade_chat failed %v", err)
		}
		author := ""
		if allchats.Data.Messages[i].Type == "trade_info" || allchats.Data.Messages[i].Type == "paxful_message" {
			author = "Villager Chief"
		} else if villagerUser == 0 {
			author = "Seller"
		} else if villagerUser == 1 {
			author = "You"
		}
		if allchats.Data.Messages[i].Type == "marked_paid" {
			UpdateTradeStatusQuery := "update trade set status = 'paid' where trade_id = ?"
			_, err = s.DB.Exec(UpdateTradeStatusQuery, tradeID)
			if err != nil {
				log.Printf("PAXFUL update trade failed because %v", err)

			}
		}

		if allchats.Data.Messages[i].Type == "released_completed" {
			UpdateTradeStatusQuery := "update trade set status = 'complete' where trade_id = ?"
			_, err = s.DB.Exec(UpdateTradeStatusQuery, tradeID)
			if err != nil {
				log.Printf("PAXFUL update trade failed status: %v", err)

			}
		}

		chatID, _ := smsObject.LastInsertId()
		chats.ID = uint64(chatID)
		chats.Text = allchats.Data.Messages[i].Text
		chats.Timestamp = uint64(allchats.Data.Messages[i].Timestamp)
		chats.Author = author

		c.WriteJSON(chats)

	}

	return newTradeResponse

}
