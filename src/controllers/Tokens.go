package controllers

import (
	"fmt"
	"log"
	"net/http"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) SupportedTokens(w http.ResponseWriter, r *http.Request) {
	// Fetch supported Stablecoins
	var res models.HttpResponse

	query := "SELECT c.crypto_currency_id,c.name,c.code FROM crypto_currency c  where c.status=1"
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Printf("Unable to fetch supported Stablecoins from DB:  %v", err)
	}
	var supportedTokens []models.SupportedTokens
	for rows.Next() {
		var tokens models.SupportedTokens
		if err = rows.Scan(&tokens.CryptoCurrencyId, &tokens.CryptoName, &tokens.CryptoCode); err != nil {
			log.Printf("unable to read supported Chain record %v", err)
		}

		supportedTokens = append(supportedTokens, tokens)
	}

	res.Data = supportedTokens
	res.Message = "Success"
	res.Status = "200"
	log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", res, res.Status))
	HttpResponse(200, res, w)

}
