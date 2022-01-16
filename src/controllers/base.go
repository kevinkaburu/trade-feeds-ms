package controllers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"sync"
	"trades/src/models"
	"trades/src/utils"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	DB      *sql.DB
	Router  *mux.Router
	RedisDB *redis.Client
}

func (s *Server) Initialize() {
	//init logger
	utils.InitLogger()
	var err error

	//init DB
	var DSN = os.Getenv("db_user") + ":" + os.Getenv("db_pass") + "@tcp(" + os.Getenv("db_host") + ":" + os.Getenv("db_port") + ")/" + os.Getenv("db_name")
	s.DB, err = sql.Open("mysql", DSN)
	if err != nil {
		log.Println("Unable to connect to db:", err)
		os.Exit(3)
	}
	log.Println("Connected to db successfully")
	s.DB.SetMaxOpenConns(100)
	s.DB.SetMaxIdleConns(64)
	s.DB.SetConnMaxIdleTime(40)

	//init Router
	s.initializeRoutes()

}
func ValidateMail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *Server) Web(wg *sync.WaitGroup) {
	defer wg.Done()
	addr := ":" + os.Getenv("API_PORT")
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "content-type", "Content-Length", "application/json", "Accept-Encoding", "Authorization", "Accept", "multipart/form-data"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("Listening to port ===>", os.Getenv("API_PORT"))
	log.Println("============ STARTING =================")

	log.Println(http.ListenAndServe(addr, handlers.CORS(originsOk, headersOk, methodsOk)(s.Router)))
}

//Health check
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"alive": true}`)
}

///Sample shared function

func (s *Server) FetchAccount(profileID int) (models.Account, error) {
	var Account models.Account
	query := "select p.profile_id,p.country_id,p.msisdn,pb.balance,pb.bonus_balance,bp.points,bp.status points_status,af.status frozen from profile p left join profile_balance pb using (profile_id) left join betika_point bp using(profile_id) left join account_freeze af using(msisdn) where profile_id=?;"
	row := s.DB.QueryRow(query, profileID)

	err := row.Scan(&Account.ProfileID, &Account.CountryID, &Account.Msisdn, &Account.Balance, &Account.BonusBalance, &Account.Points, &Account.PointsStatus, &Account.Frozen)
	if err != nil {
		return Account, err
	}
	return Account, nil

}
func IntDigitsCount(number int) int {
	count := 0
	for number != 0 {
		number /= 10
		count += 1
	}
	return count

}

///Response proccessor
func HttpResponse(statusCode int, betResponse models.HttpResponse, w http.ResponseWriter) {
	response, err := json.Marshal(betResponse)
	if err != nil {
		log.Println("Unable to parse Json BetResponse because ", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(string(response)))
	return
}

/*
Auth helper
- Based on payload submitted fetch auth data from redis.
- Validate and return True/False

*/

func AuthenticateUser(r *http.Request) bool {
	var authenticated bool
	//TO-DO
	// write auth login here.

	return authenticated
}
