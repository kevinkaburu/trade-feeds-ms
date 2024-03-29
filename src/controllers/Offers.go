package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) Offers(w http.ResponseWriter, r *http.Request) {

	var response models.HttpResponse
	statusCode := 400

	offerParamsBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to parse Offers request json, error: ", err)
		response.Message = "Unable to parse request"
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
	//const perc = (((Number(offer.fiat_price_per_crypto) - Number(offer.usd_exchange_rate)) / Number(offer.usd_exchange_rate)) * 100).toFixed(2)

	selectOfferQuery := "select ols.last_seen, ROUND((((o.fiat_price_per_crypto-fe.usd_exchange) / fe.usd_exchange)*100 ),2) diff,cc.crypto_currency_id,fe.usd_exchange, fe.forex_exchange_id, c.chain_name,c.chain_code,c.block_chain_id,cy.country_id,cy.iso_code,cy.name country_name,o.offer_id,c.block_chain_id,o.type,o.min_fiat_amount,o.max_fiat_amount,f.currency_code fiat_code,cc.code crypto_code,o.fiat_price_per_crypto,o.created,(o.max_fiat_amount/o.fiat_price_per_crypto) max_crypto from offer o inner join chain c using (chain_id) inner join fiat_currency f using (fiat_currency_id) inner join forex_exchange fe using(fiat_currency_id) inner join country cy using(country_id) inner join crypto_currency cc using (crypto_currency_id) inner join offer_payment_method opm using (offer_id) inner join payment_method pm using (payment_method_id)  left join offer_last_seen ols using (offer_id) where o.status = 1 and pm.payment_type_id != 5 %s  group by o.offer_id order by 2 asc limit 150"
	whereAppend := ""

	if int(offerParams.CountryID) > 0 {
		whereAppend += fmt.Sprintf(" and f.country_id=%d", int(offerParams.CountryID))

	}
	if int(offerParams.CryptoID) > 0 {
		whereAppend += fmt.Sprintf(" and cc.crypto_currency_id=%d", int(offerParams.CryptoID))
	}

	if int(offerParams.BlockchainID) > 0 {
		whereAppend += fmt.Sprintf(" and c.block_chain_id = %d", int(offerParams.BlockchainID))
	}

	if int(offerParams.PaymentMethodID) > 0 {
		whereAppend += fmt.Sprintf(" and opm.payment_method_id = %d", int(offerParams.PaymentMethodID))
	}

	query := fmt.Sprintf(selectOfferQuery, whereAppend)
	//fmt.Println(query)
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to query Offers from DB: %s| error: %v ", query, err))
	}

	var Offers []models.OfferDbQuery
	var lastSeen int64
	for rows.Next() {
		var offer models.OfferDbQuery
		if err = rows.Scan(&lastSeen, &offer.PercentDiff, &offer.CryptoCurrencyId, &offer.ForexExchangeUsd, &offer.ForexExchangeID, &offer.ChainName, &offer.ChainCode, &offer.ChainID, &offer.CountryID, &offer.CountryCode, &offer.CountryName, &offer.OfferID, &offer.BlockChainId, &offer.Type, &offer.MinFiatAmount, &offer.MaxFiatAmount, &offer.FiatCode, &offer.CryptoCode, &offer.FiatPricePerCrypto, &offer.Created, &offer.MaxCrypto); err != nil {
			log.Printf("unable to read Offer record %v", err)
		}

		//fetch payment modes
		PmodeQuery := fmt.Sprintf("select pm.payment_method_id,opm.tags, pm.label payment_method,pt.name payment_type from offer_payment_method opm inner join payment_method pm using (payment_method_id) inner join payment_type pt using (payment_type_id) where opm.offer_id = %v", offer.OfferID)
		paymentRows, err := s.DB.Query(PmodeQuery)
		if err != nil {
			log.Println(fmt.Sprintf("Unable to query paymentRows from DB: %s| error: %v ", PmodeQuery, err))
		}
		var paymentmodes []models.PaymentMode
		for paymentRows.Next() {
			var pmode models.PaymentMode
			if err = paymentRows.Scan(&pmode.PaymentMethodId, &pmode.Tags, &pmode.PaymentMethod, &pmode.PaymentType); err != nil {
				log.Printf("unable to read paymentOptions record %v", err)
			}
			paymentmodes = append(paymentmodes, pmode)

		}
		now := time.Now()
		timeNow := now.Unix()
		lastSeenSecs := timeNow - lastSeen

		seen := (lastSeenSecs / 60)

		offer.Seen = int(seen)

		offer.CountryFlag = fmt.Sprintf("https://flagcdn.com/32x24/%s.png", offer.CountryCode)

		offer.Payment = paymentmodes

		Offers = append(Offers, offer)

	}
	var OfferList models.OfferList
	OfferList.Count = len(Offers)
	OfferList.Offers = Offers

	response.Data = OfferList
	response.Message = "Success"

	statusCode = 200
	response.Status = fmt.Sprintf("%d", statusCode)
	//log.Println(fmt.Sprintf("Processed | StatusCode: %v ", statusCode))
	HttpResponse(statusCode, response, w)

}
