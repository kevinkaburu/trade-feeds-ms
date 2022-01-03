package controllers

import (
	"fmt"
	"log"
	"net/http"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	// Process significant events.

	var res models.HttpResponse

	statusCode, response := 200, res
	log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", response, statusCode))
	HttpResponse(statusCode, response, w)
	return

}
