package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) Kyc(w http.ResponseWriter, r *http.Request) {
	// Process significant events.
	bytesValue, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to parse ValidateBet request json, error: ", err, bytesValue)
		return
	}

	var res models.HttpResponse

	statusCode, response := 200, res
	log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", response, statusCode))
	HttpResponse(statusCode, response, w)
	return

}
