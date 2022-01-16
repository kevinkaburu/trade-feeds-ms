package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) Offers(w http.ResponseWriter, r *http.Request) {

	var response models.HttpResponse
	statusCode := 400

	offerParamsBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to parse Signup request json, error: ", err)
		response.Message = "Unable to parse Signup request"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	var offerParams models.OfferQuery

	if len(offerParamsBytes) > 0 {

		err = json.Unmarshal(offerParamsBytes, &offerParams)
		if err != nil {
			log.Println("Unable to parse  OfferQuery Json  because ", err)
			response.Message = "Unable to parse  OfferQuery Json"
			response.Status = "400"
			HttpResponse(statusCode, response, w)
			return
		}
	}
	/*
			CountryID    StringInt   `json:"country_id"`
			Fiat         string      `json:"Fiat"`
			FiatAmount   StringFloat `json:"fiat_amount"`
			CryptoAmount StringFloat `json:"crypto_amount"`


			//
		OfferID            int     `json:"offer_id"`
		Type               string  `json:"type"`
		MinFiatAmount      float64 `json:"min_fiat_amount"`
		MaxFiatAmount      float64 `json:"max_fiat_amount"`
		FiatCode           string  `json:"fiat_code"`
		CryptoCode         string  `json:"crypto_code"`
		FiatPricePerCrypto float64 `json:"fiat_price_per_crypto"`
		Created            string  `json:"created"`


	*/
	selectOfferQuery := "select o.offer_id,o.type,o.min_fiat_amount,o.max_fiat_amount,f.currency_code fiat_code,cc.code crypto_code,o.fiat_price_per_crypto,o.created,(o.max_fiat_amount/o.fiat_price_per_crypto) max_crypto from offer o inner join fiat_currency f using (fiat_currency_id) inner join crypto_currency cc using (crypto_currency_id) where o.status = 1 %s order by o.fiat_price_per_crypto asc limit 100"
	whereAppend := ""

	if int(offerParams.CountryID) > 0 {
		whereAppend += fmt.Sprintf(" and f.country_id=%d", int(offerParams.CountryID))

	}
	if len(offerParams.Fiat) == 3 {
		whereAppend += fmt.Sprintf(" and f.currency_code='%s'", offerParams.Fiat)
	}

	if float64(offerParams.FiatAmount) > 0.0 {
		whereAppend += fmt.Sprintf(" and o.max_fiat_amount > %f", float64(offerParams.FiatAmount))
	}

	if float64(offerParams.CryptoAmount) > 0.0 {
		whereAppend += fmt.Sprintf(" and crptos > %f", float64(offerParams.CryptoAmount))
	}

	query := fmt.Sprintf(selectOfferQuery, whereAppend)
	fmt.Println(query)
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to query Offers from DB: %s| error: %v ", query, err))
	}
	var Offers []models.OfferDbQuery
	for rows.Next() {
		var offer models.OfferDbQuery
		if err = rows.Scan(&offer.OfferID, &offer.Type, &offer.MinFiatAmount, &offer.MaxFiatAmount, &offer.FiatCode, &offer.CryptoCode, &offer.FiatPricePerCrypto, &offer.Created, &offer.MaxCrypto); err != nil {
			log.Printf("unable to read Offer record %v", err)
		}
		Offers = append(Offers, offer)

	}
	fmt.Println(Offers)

	response.Data = Offers

	statusCode, response = 200, response
	log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", response, statusCode))
	HttpResponse(statusCode, response, w)
	return

}
