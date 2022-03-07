package controllers

import (
	"github.com/gorilla/mux"
)

func (s *Server) initializeRoutes() {
	s.Router = mux.NewRouter()
	//health check
	s.Router.HandleFunc("/", HealthCheckHandler).Methods("GET")
	// User MS endpoints
	s.Router.HandleFunc("/trade/start", s.StartTrade).Methods("POST")
	s.Router.HandleFunc("/trade/chat", s.TradeChat)
	s.Router.HandleFunc("/user/login", s.Login).Methods("POST")
	s.Router.HandleFunc("/trade/offer", s.Offers).Methods("POST")
	s.Router.HandleFunc("/user/logout", s.Logout).Methods("POST")
	s.Router.HandleFunc("/trade/chains", s.SupportedChains).Methods("GET")
	s.Router.HandleFunc("/trade/countries", s.SupportedCountries).Methods("GET")
	s.Router.HandleFunc("/trade/paymentmethods", s.SupportedPaymentMethods).Methods("GET")
	s.Router.HandleFunc("/trade/Stablecoins", s.SupportedTokens).Methods("GET")

}
