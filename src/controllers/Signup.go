package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"trades/src/models"
	"trades/src/utils"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	// Process significant events.
	var response models.HttpResponse
	var statusCode int
	statusCode = 400
	signupBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to parse Signup request json, error: ", err)
		response.Message = "Unable to parse Signup request"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	var signuppayload models.SignupPayload

	if err = json.Unmarshal(signupBytes, &signuppayload); err != nil {
		log.Println("Unable to parse  Signup Json  because ", err)
		response.Message = "Unable to parse  Signup Json"
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	//Validate inputs
	//If only username check if it exists in DB
	//validate email
	if signuppayload.Email == "" && signuppayload.Msisdn == 0 {
		log.Println("Email or mobile required")
		response.Message = "Email or mobile required."
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	//if Mail in the payload
	if len(signuppayload.Email) > 1 {
		if !ValidateMail(signuppayload.Email) {
			log.Println("Email not valid ", signuppayload)
			response.Message = "Invalid email address."
			response.Status = "400"
			HttpResponse(statusCode, response, w)
			return
		}
		if !strings.Contains(signuppayload.Email, ".") {
			log.Println("Email not valid ", signuppayload)
			response.Message = "Invalid email address."
			response.Status = "400"
			HttpResponse(statusCode, response, w)
			return
		}

		//check if mail already exists.
		var profileID string
		query := "select profile_id from profile where email=?;"
		row := s.DB.QueryRow(query, signuppayload.Email)

		_ = row.Scan(&profileID)
		if profileID != "" {
			log.Println("Email Exists. ProfileID: ", profileID)
			response.Message = "Email address already exists."
			response.Status = "400"
			HttpResponse(statusCode, response, w)
			return
		}
	}

	if signuppayload.Msisdn > 1 {
		msisdnLength := IntDigitsCount(int(signuppayload.Msisdn))
		if msisdnLength < 11 {
			log.Println("Malformed msisdn with length of: ", msisdnLength)
			response.Message = "mobile number should include country code."
			response.Status = "400"
			HttpResponse(statusCode, response, w)
			return
		}

		var profileID string
		query := "select profile_id from profile where msisdn=?;"
		row := s.DB.QueryRow(query, signuppayload.Msisdn)

		_ = row.Scan(&profileID)
		if profileID != "" {
			log.Println("Mobile number Exists. ProfileID: ", profileID)
			response.Message = "Mobile number already exists."
			response.Status = "400"
			HttpResponse(statusCode, response, w)
			return
		}

	}
	//passowrd validation
	if signuppayload.ConfirmPassword != signuppayload.Password {
		log.Println("Confirm Pass and Password are not same.")
		response.Message = "Password and confirmation password do not match."
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}

	passwordValidationError := VerifyPassword(signuppayload.Password)
	if passwordValidationError != nil {
		log.Println("Password ", passwordValidationError)
		response.Message = passwordValidationError.Error()
		response.Status = "400"
		HttpResponse(statusCode, response, w)
		return
	}
	log.Println("Ready to save to DB", signuppayload)
	//Hashpassword
	HashedPassword, error := HashPassword(signuppayload.Password)
	if error != nil {
		log.Println("Failed to hash password")
		response.Message = "Internal server Error."
		response.Status = "500"
		HttpResponse(statusCode, response, w)
		return

	}

	//Save payload to DB
	stmt, err := s.DB.Prepare("insert into profile set email=?, mobile= ? , password=?,status=?,created=now(),modified=now()")
	if err != nil {
		log.Println("Failed on query prepare (insert into profile)")
		response.Message = "Internal server Error."
		response.Status = "500"
		HttpResponse(statusCode, response, w)
		return
	}
	defer stmt.Close()

	// execute
	res, err := stmt.Exec(signuppayload.Email, signuppayload.Msisdn, HashedPassword, 1)
	if err != nil {
		log.Println("Failed ON statement Exec")
		response.Message = "Internal server Error."
		response.Status = "500"
		HttpResponse(statusCode, response, w)
		return
	}
	profileID, err := res.LastInsertId()
	if err != nil {
		log.Println("Failed to get LastInsertId")
		response.Message = "Internal server Error."
		response.Status = "500"
		HttpResponse(statusCode, response, w)
		return
	}
	log.Println(fmt.Sprintf("ProfileID response: %v", profileID))

	//Send opt/ mail verification
	go utils.SendConfirmMail(signuppayload.Email, fmt.Sprintf("%v", profileID))

	//Auth processor (to redis/ encrypt)

	//return

	statusCode = 201
	log.Println(fmt.Sprintf("Process response: %s| StatusCode: %v ", response, statusCode))
	HttpResponse(statusCode, response, w)
	return

}
