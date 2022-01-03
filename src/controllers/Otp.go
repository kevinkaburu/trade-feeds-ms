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

func (s *Server) Otp(w http.ResponseWriter, r *http.Request) {
	// Process significant events.
	bytesValue, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to parse ValidateBet request json, error: ", err)
		return
	}

	var betData models.Account
	//log.Println(fmt.Sprintf("Successful request received payload: %v", string(bytesValue)))

	if err = json.Unmarshal(bytesValue, &betData); err != nil {
		log.Println("Unable to parse Json Bet because ", err)
		return
	}

	var res models.HttpResponse

	statusCode, response := 200, res
	log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", response, statusCode))
	HttpResponse(statusCode, response, w)
	return

}
