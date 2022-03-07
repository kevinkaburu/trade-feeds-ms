package controllers

import (
	"fmt"
	"log"
	"net/http"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) SupportedPaymentMethods(w http.ResponseWriter, r *http.Request) {
	// Fetch supported Payment methods
	var res models.HttpResponse

	query := "SELECT payment_method_id,label,payment_type_id FROM `payment_method` order by label asc"
	//fmt.Println(query)
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Printf("Unable to fetch supported paymentMethods from DB:  %v", err)
	}
	var supportedPaymentMethods []models.SupportedPaymentMethods
	for rows.Next() {
		var pmethods models.SupportedPaymentMethods
		if err = rows.Scan(&pmethods.PaymentMethodId, &pmethods.PaymentMethodName, &pmethods.PaymentTypeID); err != nil {
			log.Printf("unable to read supported Chain record %v", err)
		}

		supportedPaymentMethods = append(supportedPaymentMethods, pmethods)
	}

	res.Data = supportedPaymentMethods
	res.Message = "Success"
	res.Status = "200"
	log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", res, res.Status))
	HttpResponse(200, res, w)

}
