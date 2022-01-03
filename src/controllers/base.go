package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"sync"
	"trades/src/models"
	"trades/src/utils"
	"unicode"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

//Validate Passowrd
func VerifyPassword(password string) error {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	const minPassLength = 6
	const maxPassLength = 64
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("atleast one numeric character required")
	}
	if !specialCharPresent {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength))
	}

	if len(errorString) != 0 {
		return fmt.Errorf("Password error: " + errorString)
	}
	return nil
}

//Password hashing
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//Password validation
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
