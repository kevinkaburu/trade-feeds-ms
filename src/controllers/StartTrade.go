package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"trades/src/models"
	"trades/src/utils"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) StartTrade(w http.ResponseWriter, r *http.Request) {
	// Process significant events.
	var response models.HttpResponse
	var statusCode int
	statusCode = 400

	tokenCookie, err := r.Cookie("vlg")
	if err != nil {
		log.Println("Unable to fetch tokenCookie error: ", err)
		response.Message = "Unable to parse request"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	//fetch data from redis/ for auth
	redisData, errr := utils.FetchDataFromRedis(tokenCookie.Value, s.RedisDB)
	if errr != nil {
		fmt.Printf("Redis Fetch Naounce Failed: %s\n", errr)
	}
	var redisWalletData models.WalletAuth
	if err = json.Unmarshal([]byte(redisData), &redisWalletData); err != nil {
		log.Println("Unable to parse  Redis  Json  because ", err)
		return
	}

	requestBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to parse StartTrade request json, error: ", err)
		response.Message = "Unable to parse request"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	var startTradePayload models.StartTradeQuery

	if err = json.Unmarshal(requestBytes, &startTradePayload); err != nil {
		log.Println("Unable to parse  StartTrade Json  because ", err)
		response.Message = "Unable to parse Json"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	//Validate inputs
	if startTradePayload.FiatAmount < 0.1 {
		log.Println("Fiat amount  required")
		response.Message = "Fiat amount  required"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	//validate Offer
	if startTradePayload.OfferID == 0 {
		log.Println("Offer  required")
		response.Message = "Offer  required"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	//Wallet account exists
	if len(startTradePayload.WalletAccount) < 26 {
		log.Println(fmt.Printf("Invalid wallet Addres: %v | Length: %v", startTradePayload.WalletAccount, len(startTradePayload.WalletAccount)))
		response.Message = "Invalid wallet Address"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	//Wallet == yo yhe current Auth'd
	if startTradePayload.WalletAccount != redisWalletData.Walletdata.Address {
		log.Println(fmt.Printf("Auth'd wallet Addres: %v | UserWallet: %v", redisWalletData.Walletdata.Address, startTradePayload.WalletAccount))
		response.Message = "Authentication failed. Connect your wallet to proceed."
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	var (
		OfferID       uint64
		ProviderID    uint64
		ProfileID     uint64
		ExternalID    string
		Status        uint64
		MinFiat       float64
		MaxFiat       float64
		PaymentMethod string
	)
	//fetch Offer from DB
	selectOfferQuery := fmt.Sprintf("select o.offer_id,o.min_fiat_amount,o.max_fiat_amount,o.provider_id,o.profile_id,o.external_id,o.status,pm.label from offer o inner join offer_payment_method opm using (offer_id) inner join payment_method pm using (payment_method_id)  where offer_id= %v", startTradePayload.OfferID)
	err = s.DB.QueryRow(selectOfferQuery).Scan(&OfferID, &MinFiat, &MaxFiat, &ProviderID, &ProfileID, &ExternalID, &Status, &PaymentMethod)
	if err != nil {
		log.Println("Unknown Offer  required")
		response.Message = "Offer  not found."
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	if float64(startTradePayload.FiatAmount) < MinFiat {
		log.Println(fmt.Printf(" %v is less than Minimum Fiat amount of %v ", float64(startTradePayload.FiatAmount), MinFiat))
		response.Message = fmt.Sprintf(" %v is less than Minimum Fiat amount of %v ", float64(startTradePayload.FiatAmount), MinFiat)
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	if float64(startTradePayload.FiatAmount) > MaxFiat {
		log.Println(fmt.Printf(" %v is greater than maximum Fiat amount of %v ", float64(startTradePayload.FiatAmount), MaxFiat))
		response.Message = fmt.Sprintf(" %v is greater than maximum Fiat amount of %v ", float64(startTradePayload.FiatAmount), MaxFiat)
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	//Confirm only paxful is live
	provider_id, err := strconv.Atoi(os.Getenv("PAXFUL_PROVIDER_ID"))
	if err != nil {
		log.Println("Unable to get Paxful Provider_id from .env")
		response.Message = "Internal Error. "
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	if ProviderID != uint64(provider_id) {
		log.Println("Only Paxful that is active at the moment")
		response.Message = "We are not able to process this Offer at the moment."
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	//Fetch the offer from Paxful
	paxfulOfferData := s.GetOffer(ExternalID)
	//if paxful status no success update db status and return error.
	if paxfulOfferData.Status != "success" {
		UpdateOffersStatusQuery := "update offer set status = 0 where offer_id = ?"
		_, err = s.DB.Exec(UpdateOffersStatusQuery, OfferID)
		if err != nil {
			log.Printf("PAXFUL get trade statu: %v | unable to update  to offer  because %v", paxfulOfferData.Status, err)

		}
		log.Println("PAXFUL get trade status non success")
		response.Message = "We are not able to process this Offer at the moment."
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return

	}

	//check if the offer is active on paxful end
	if !paxfulOfferData.Data.Active {
		UpdateOffersStatusQuery := "update offer set status = 0 where offer_id = ?"
		_, err = s.DB.Exec(UpdateOffersStatusQuery, OfferID)
		if err != nil {
			log.Printf("PAXFUL get trade Active: %v | unable to update  to offer  because %v", paxfulOfferData.Data.Active, err)

		}
		log.Println("PAXFUL Trade not active.")
		response.Message = "We are not able to process this Offer at the moment."
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	//create  Trade in DB
	inserTrade := "insert into trade set profile_id=?, offer_id=?,fiat_amount=?,status=?,wallet=?,created=now(),modified=now()"
	tradeObject, err := s.DB.Exec(inserTrade, ProfileID, OfferID, startTradePayload.FiatAmount, "started", startTradePayload.WalletAccount)
	if err != nil {
		log.Printf("Create db trade error: %v", err)

	}
	tradeID, err := tradeObject.LastInsertId()
	if err != nil {
		log.Printf("Unable to get TradeID error: %v", err)
	}
	//Fanya ile kitu. Call paxful and init a trade
	log.Println(fmt.Sprintf("TradeID created ... : %d", tradeID))
	NewTrade := s.CreateTrade(paxfulOfferData.Data.OfferHash, float64(startTradePayload.FiatAmount))

	if NewTrade.Status != "success" {
		log.Printf(fmt.Sprintf("Unalbe it init trade. Response: %v", NewTrade))
		UpdateTradeStatusQuery := "update trade set status = 'cancelled' where trade_id = ?"
		_, err = s.DB.Exec(UpdateTradeStatusQuery, tradeID)
		if err != nil {
			log.Printf("PAXFUL update trade failed status: %v | unable to update  to offer  because %v", NewTrade.Status, err)

		}
		log.Printf(" Unlable to trade: %v", NewTrade)
		response.Message = "We are not able to process this Offer at the moment."
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return

	}
	///it's created
	UpdateTradeStatusQuery := "update trade set external_id = ?  where trade_id = ?"
	_, err = s.DB.Exec(UpdateTradeStatusQuery, NewTrade.Data.TradeHash, tradeID)
	if err != nil {
		log.Printf("PAXFUL update trade failed status: %v | unable to update  to offer  because %v", NewTrade.Status, err)
	}

	//Send Message Seller requesting for payment details
	sms := fmt.Sprintf("Hello,\n Share your <b> %v </b> payment details.", PaymentMethod)
	messageSent := s.SendMessage(NewTrade.Data.TradeHash, sms)
	if !messageSent {
		//cancel transaction and revert
		//for now we just log
		log.Printf("PAXFUL Unable to send message")
	}

	//Fetch chats.
	InitialMessages := s.FetchMessages(NewTrade.Data.TradeHash)
	//init response to feed it with the messages
	var newTradeResponse models.StartTradeResponse
	//loop through them inserting into db.
	for i := 0; i < len(InitialMessages.Data.Messages); i++ {
		var chats models.ChatMessages
		InsertChatQuery := "insert ignore into trade_chat set trade_id=?,type=?,external_id=?,author=?,text=?,villager_user=?,timestamp=?,created=now(),modified=now()"
		villagerUser := 0
		if InitialMessages.Data.Messages[i].Author == os.Getenv("VILLAGER_PAXFUL_USERNAME") {
			villagerUser = 1
		}
		smsObject, err := s.DB.Exec(InsertChatQuery, tradeID, InitialMessages.Data.Messages[i].Type, InitialMessages.Data.Messages[i].ID, InitialMessages.Data.Messages[i].Author, InitialMessages.Data.Messages[i].Text, villagerUser, InitialMessages.Data.Messages[i].Timestamp)
		if err != nil {
			log.Printf("PAXFUL insert ignore trade_chat failed %v", err)
		}
		author := ""
		if InitialMessages.Data.Messages[i].Type == "trade_info" || InitialMessages.Data.Messages[i].Type == "paxful_message" {
			author = "Villager Chief"
		} else if villagerUser == 0 {
			author = "Seller"
		} else if villagerUser == 1 {
			author = "You"
		}
		chatID, _ := smsObject.LastInsertId()
		chats.ID = uint64(chatID)
		chats.Text = InitialMessages.Data.Messages[i].Text
		chats.Timestamp = uint64(InitialMessages.Data.Messages[i].Timestamp)
		chats.Author = author

		newTradeResponse.Messages = append(newTradeResponse.Messages, chats)

	}

	//construct response

	newTradeResponse.TradeID = tradeID
	newTradeResponse.TradeStatus = "STARTED"
	newTradeResponse.PaymentMethod = PaymentMethod

	//fetch the messages to get their IDs

	response.Data = newTradeResponse

	statusCode = 201
	//log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", response, statusCode))
	HttpResponse(statusCode, response, w)

}

func (s *Server) GetOffer(offerID string) (offerdata models.OfferGetData) {
	data := url.Values{}
	data.Set("offer_hash", offerID)
	endpoint := fmt.Sprintf("%s/offer/get", os.Getenv("PAXFUL_BASE_URL"))
	//http request
	resp, err := s.PaxfulClient.PostForm(endpoint, data)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: %v", err)

	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 202 {
		if err = json.Unmarshal(body, &offerdata); err != nil {
			fmt.Print("Unable to read response into struct because ", err)

		}
		return

	}
	return
}

func (s *Server) CreateTrade(offerID string, fiatAmount float64) (tradeStart models.PaxfulTradeStart) {
	data := url.Values{}
	data.Set("offer_hash", offerID)
	data.Set("fiat", fmt.Sprintf("%f", fiatAmount))
	//endpoint := fmt.Sprintf("%s/trade/start", os.Getenv("PAXFUL_BASE_URL"))
	endpoint := ""
	//http request
	resp, err := s.PaxfulClient.PostForm(endpoint, data)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: %v", err)
		return

	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 202 {
		if err = json.Unmarshal(body, &tradeStart); err != nil {
			fmt.Print("Unable to read response into struct because ", err)
		}
		return

	}
	return
}
