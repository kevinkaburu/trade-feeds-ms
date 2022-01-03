package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
)

func (s *Server) PaxfulFetchOffers(wg *sync.WaitGroup) {
	defer wg.Done()
	// Fetch active offers from paxful

	pticker := time.NewTicker(10 * time.Second)
	pquit := make(chan struct{})
	go func() {
		for {
			select {
			case <-pticker.C:
				s.FetchActiveOffers()

			case <-pquit:
				pticker.Stop()
				return
			}
		}
	}()
}

func (s *Server) FetchActiveOffers() {

	//Get token
	config := clientcredentials.Config{
		ClientID:     os.Getenv("PAXFUL_VILLAGERS_APP_ID"),
		ClientSecret: os.Getenv("PAXFUL_VILLAGERS_SECRET"),
		TokenURL:     os.Getenv("PAXFUL_ACCESS_TOKEN_URL"),
		Scopes:       []string{},
	}
	//setup context
	client := config.Client(context.Background())
	//fetch possible currencies from db

	activeFiatsQuery := "select fiat_currency_id,currency_code from fiat_currency where status=1"
	rows, err := s.DB.Query(activeFiatsQuery)
	if err != nil {
		fmt.Printf("unable to fetch fiat Currencies: %v| error: %v", activeFiatsQuery, err)
		return
	}

	for rows.Next() {
		var fiats models.FiatCurency
		if err = rows.Scan(&fiats.FiatCurrencyId, &fiats.FiatCurrencyCode); err != nil {
			log.Printf("unable to read Currency record %v", err)
			continue
		}
		go s.httpOffers(client, fiats)
	}

	return
}

func (s *Server) httpOffers(c *http.Client, fiatCurrency models.FiatCurency) {
	data := url.Values{}
	data.Set("offer_type", "buy")
	data.Set("type", "buy")
	data.Set("limit", "200")
	data.Set("currency_code", fiatCurrency.FiatCurrencyCode)
	data.Set("crypto_currency_code", "usdt")
	endpoint := fmt.Sprintf("%s/offer/all", os.Getenv("PAXFUL_BASE_URL"))
	//http request
	resp, err := c.PostForm(endpoint, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)

	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 202 {
		var paxfulOffers models.PaxfulOffers

		if err = json.Unmarshal(body, &paxfulOffers); err != nil {
			fmt.Print("Unable to read response into struct because ", err)

		}
		for i := 0; i < len(paxfulOffers.Data.Offers); i++ {
			provider_id := os.Getenv("PAXFUL_PROVIDER_ID")
			profile_id := os.Getenv("VILLAGER_PROFILE_ID")
			crypto_currency_id := os.Getenv("CRYPTO_USDT_ID")
			offertype := paxfulOffers.Data.Offers[i].OfferType
			external_id := paxfulOffers.Data.Offers[i].OfferID
			min_fiat_amount := paxfulOffers.Data.Offers[i].FiatAmountRangeMin
			max_fiat_amount := paxfulOffers.Data.Offers[i].FiatAmountRangeMax
			fiat_price_per_crypto := paxfulOffers.Data.Offers[i].FiatPricePerCrypto

			insertOffer := "insert  into offer (provider_id,profile_id,type,external_id,min_fiat_amount,max_fiat_amount,fiat_currency_id,crypto_currency_id,fiat_price_per_crypto,created,modified) VALUES (?,?,?,?,?,?,?,?,?,now(),now()) ON DUPLICATE KEY UPDATE min_fiat_amount=?,max_fiat_amount=?,fiat_price_per_crypto=?, modified = now();"
			if _, err := s.DB.Exec(insertOffer, provider_id, profile_id, offertype, external_id, min_fiat_amount, max_fiat_amount, fiatCurrency.FiatCurrencyId, crypto_currency_id, fiat_price_per_crypto, min_fiat_amount, max_fiat_amount, fiat_price_per_crypto); err != nil {
				log.Printf("unable to insert to db to %v because %v", paxfulOffers.Data.Offers[i].OfferID, err)
				return
			}

		}

	} else {
		fmt.Print("Unable to fetch offers", string(body))
		err = errors.New("Error fetching offers")
	}
}
