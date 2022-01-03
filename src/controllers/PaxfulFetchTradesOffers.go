package controllers

import (
	"fmt"
	"log"
	"sync"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) PaxfulFetchOffers(wg *sync.WaitGroup) {
	// Fetch active offers from paxful

	var res models.HttpResponse

	statusCode, response := 200, res
	log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", response, statusCode))
	return

}
