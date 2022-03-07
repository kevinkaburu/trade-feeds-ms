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

	selectOfferQuery := "select o.offer_id,c.block_chain_id,o.type,o.min_fiat_amount,o.max_fiat_amount,f.currency_code fiat_code,cc.code crypto_code,o.fiat_price_per_crypto,o.created,(o.max_fiat_amount/o.fiat_price_per_crypto) max_crypto from offer o inner join chain c using (chain_id) inner join fiat_currency f using (fiat_currency_id) inner join crypto_currency cc using (crypto_currency_id) inner join offer_payment_method opm using (offer_id) inner join payment_method pm using (payment_method_id)  where o.status = 1 and pm.payment_type_id != 5 %s  group by o.offer_id order by o.offer_id,o.fiat_price_per_crypto asc limit 150"
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
	//fmt.Println(query)
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to query Offers from DB: %s| error: %v ", query, err))
	}
	var Offers []models.OfferDbQuery
	for rows.Next() {
		var offer models.OfferDbQuery
		if err = rows.Scan(&offer.OfferID, &offer.BlockChainId, &offer.Type, &offer.MinFiatAmount, &offer.MaxFiatAmount, &offer.FiatCode, &offer.CryptoCode, &offer.FiatPricePerCrypto, &offer.Created, &offer.MaxCrypto); err != nil {
			log.Printf("unable to read Offer record %v", err)
		}

		//fetch payment modes
		PmodeQuery := fmt.Sprintf("select opm.tags, pm.label payment_method,pt.name payment_type from offer_payment_method opm inner join payment_method pm using (payment_method_id) inner join payment_type pt using (payment_type_id) where opm.offer_id = %v", offer.OfferID)
		paymentRows, err := s.DB.Query(PmodeQuery)
		if err != nil {
			log.Println(fmt.Sprintf("Unable to query paymentRows from DB: %s| error: %v ", PmodeQuery, err))
		}
		var paymentmodes []models.PaymentMode
		for paymentRows.Next() {
			var pmode models.PaymentMode
			if err = paymentRows.Scan(&pmode.Tags, &pmode.PaymentMethod, &pmode.PaymentType); err != nil {
				log.Printf("unable to read paymentOptions record %v", err)
			}
			paymentmodes = append(paymentmodes, pmode)

		}
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
