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

}
