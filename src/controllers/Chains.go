package controllers

import (
	"log"
	"net/http"
	"trades/src/models"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) SupportedChains(w http.ResponseWriter, r *http.Request) {
	// Fetch supported chains
	var res models.HttpResponse

	query := "SELECT chain_id,chain_name,block_chain_id,chain_code FROM `chain` where status= 1 order by chain_id asc"
	//fmt.Println(query)
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Printf("Unable to fetch supported Chains from DB:  %v", err)
	}
	var supportedChains []models.SupportedChains
	var LchainID int = 0
	for rows.Next() {
		var chains models.SupportedChains
		if err = rows.Scan(&LchainID, &chains.ChainName, &chains.BlockChainId, &chains.ChainCode); err != nil {
			log.Printf("unable to read supported Chain record %v", err)
		}

		supportedChains = append(supportedChains, chains)
	}

	res.Data = supportedChains
	res.Message = "Success"
	res.Status = "200"
	HttpResponse(200, res, w)

}
