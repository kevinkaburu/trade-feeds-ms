package controllers

import (
	"fmt"
	"log"
	"net/http"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) SupportedCountries(w http.ResponseWriter, r *http.Request) {
	// Fetch supported Countries
	var res models.HttpResponse

	query := "SELECT c.country_id,c.name,c.iso_code,fc.fiat_currency_id,fc.currency_code FROM country c inner join fiat_currency fc using (country_id) where fc.status=1 order by c.name asc "
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Printf("Unable to fetch supported Chains from DB:  %v", err)
	}
	var supportedCountries []models.SupportedCountries
	for rows.Next() {
		var countries models.SupportedCountries
		if err = rows.Scan(&countries.CountryID, &countries.CountryName, &countries.CountryCode, &countries.FiatCurrencyId, &countries.FiatCurencyCode); err != nil {
			log.Printf("unable to read supported Chain record %v", err)
		}
		countries.CountryFlag = fmt.Sprintf("https://flagcdn.com/16x12/%s.png", countries.CountryCode)

		supportedCountries = append(supportedCountries, countries)
	}

	res.Data = supportedCountries
	res.Message = "Success"
	res.Status = "200"
	HttpResponse(200, res, w)

}
