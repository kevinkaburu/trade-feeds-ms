package TradesMs

import (
	"sync"
	"trades/src/controllers"
)

var server = controllers.Server{}

//initialise server

func Run() {
	server.Initialize()
	var wg sync.WaitGroup

	wg.Add(2)
	go server.Web(&wg)
	go server.PaxfulFetchOffers(&wg)
	wg.Wait()

}
