package controllers

import (
	"github.com/gorilla/mux"
)

func (s *Server) initializeRoutes() {
	s.Router = mux.NewRouter()
	//health check
	s.Router.HandleFunc("/", HealthCheckHandler).Methods("GET")
	// User MS endpoints
	s.Router.HandleFunc("/user/signup", s.Signup).Methods("POST")
	s.Router.HandleFunc("/user/login", s.Login).Methods("POST")
	s.Router.HandleFunc("/user/verify/email", s.VerifyMail).Methods("GET")
	s.Router.HandleFunc("/user/verify/otp", s.Otp).Methods("POST")
	s.Router.HandleFunc("/user/kyc", s.Kyc).Methods("POST")
	s.Router.HandleFunc("/user/logout", s.Logout).Methods("POST")

}
