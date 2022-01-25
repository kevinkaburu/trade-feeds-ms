package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"sync"
	"trades/src/models"
	"trades/src/utils"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
)

type Server struct {
	DB           *sql.DB
	Router       *mux.Router
	RedisDB      *redis.Client
	PaxfulClient *http.Client
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
	//init Http CLient
	//Get token
	config := clientcredentials.Config{
		ClientID:     os.Getenv("PAXFUL_VILLAGERS_APP_ID"),
		ClientSecret: os.Getenv("PAXFUL_VILLAGERS_SECRET"),
		TokenURL:     os.Getenv("PAXFUL_ACCESS_TOKEN_URL"),
		Scopes:       []string{},
	}
	//setup context
	s.PaxfulClient = config.Client(context.Background())

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

func (s *Server) FetchMessages(tradeID string) (tradechats models.TradeChat) {
	data := url.Values{}
	data.Set("trade_hash", tradeID)
	endpoint := fmt.Sprintf("%s/trade-chat/get", os.Getenv("PAXFUL_BASE_URL"))
	//http request
	resp, err := s.PaxfulClient.PostForm(endpoint, data)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: %v", err)

	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 202 {
		if err = json.Unmarshal(body, &tradechats); err != nil {
			fmt.Print("Unable to read response into struct because ", err)
		}
		return

	}
	return
}

func (s *Server) SendMessage(tradeID, sms string) (status bool) {
	var newChatResponse models.NewChatResponse
	status = false
	data := url.Values{}
	data.Set("trade_hash", tradeID)
	data.Set("message", sms)
	endpoint := fmt.Sprintf("%s/trade-chat/post", os.Getenv("PAXFUL_BASE_URL"))
	//http request
	resp, err := s.PaxfulClient.PostForm(endpoint, data)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: %v", err)

	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 202 {
		if err = json.Unmarshal(body, &newChatResponse); err != nil {
			fmt.Print("Unable to read response into struct because ", err)
			return
		}
		if newChatResponse.Status == "success" {
			status = true
		}
		return

	}
	return
}
