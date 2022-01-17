package models

import (
	"database/sql"
	"encoding/json"
	"strconv"
)

type SignupPayload struct {
	Email           string    `json:"email"`
	Msisdn          StringInt `json:"msisdn"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
}

type OfferDbQuery struct {
	OfferID            int           `json:"offer_id"`
	Type               string        `json:"type"`
	MinFiatAmount      float64       `json:"min_fiat_amount"`
	MaxFiatAmount      float64       `json:"max_fiat_amount"`
	FiatCode           string        `json:"fiat_code"`
	CryptoCode         string        `json:"crypto_code"`
	Payment            []PaymentMode `json:"payment"`
	FiatPricePerCrypto float64       `json:"fiat_price_per_crypto"`
	Created            string        `json:"created"`
	MaxCrypto          float64       `json:"max_crypto"`
}

type PaymentMode struct {
	Tags          string `json:"tags"`
	PaymentMethod string `json:"payment_method"`
	PaymentType   string `json:"payment_type"`
}

type OfferQuery struct {
	CountryID    StringInt   `json:"country_id"`
	Fiat         string      `json:"fiat_code"`
	FiatAmount   StringFloat `json:"fiat_amount"`
	CryptoAmount StringFloat `json:"crypto_amount"`
}

type LoginPayload struct {
	Email    string    `json:"email"`
	Msisdn   StringInt `json:"msisdn"`
	Password string    `json:"password"`
}

type OtpPayload struct {
	Email  string    `json:"email"`
	Msisdn StringInt `json:"msisdn"`
	Otp    StringInt `json:"otp"`
}

type Account struct {
	ProfileID    int
	CountryID    int
	Msisdn       int
	Balance      sql.NullFloat64
	BonusBalance sql.NullFloat64
	Points       sql.NullFloat64
	PointsStatus sql.NullString
	Frozen       sql.NullInt64
}
type HttpResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

type StringInt int
type StringFloat float64

func (st *StringInt) UnmarshalJSON(b []byte) error {
	//convert the bytes into an interface
	//this will help us check the type of our value
	//if it is a string that can be converted into an int we convert it
	///otherwise we return an error
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int(v))
	case string:
		///here convert the string into
		///an integer
		if v == "" {
			v = "0"
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			///the string might not be of integer type
			///so return an error
			return err

		}
		*st = StringInt(i)

	}
	return nil
}

func (st *StringFloat) UnmarshalJSON(b []byte) error {
	//convert the bytes into an interface
	//this will help us check the type of our value
	//if it is a string that can be converted into an int we convert it
	///otherwise we return an error
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringFloat(float64(v))
	case float64:
		*st = StringFloat(v)
	case string:
		///here convert the string into
		///an integer
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			///the string might not be of integer type
			///so return an error
			return err

		}
		*st = StringFloat(i)

	}
	return nil
}

type ForexExchange struct {
	Status    string `json:"status"`
	Timestamp int    `json:"timestamp"`
	Data      struct {
		Count      int `json:"count"`
		Currencies []struct {
			Code              string `json:"code"`
			Name              string `json:"name"`
			NameLocalized     string `json:"name_localized"`
			MinTradeAmountUsd string `json:"min_trade_amount_usd"`
			Rate              struct {
				Usd  float64 `json:"usd"`
				Btc  float64 `json:"btc"`
				Usdt float64 `json:"usdt"`
				Eth  float64 `json:"eth"`
			} `json:"rate"`
		} `json:"currencies"`
	} `json:"data"`
}

type FiatCurency struct {
	FiatCurrencyId   int
	FiatCurrencyCode string
}

//Paxful  - Offer/all
type PaxfulOffers struct {
	Status    string `json:"status"`
	Timestamp int    `json:"timestamp"`
	Data      struct {
		Limit      int `json:"limit"`
		Offset     int `json:"offset"`
		Count      int `json:"count"`
		Totalcount int `json:"totalCount"`
		Offers     []struct {
			OfferID                    string      `json:"offer_id"`
			OfferType                  string      `json:"offer_type"`
			CurrencyCode               string      `json:"currency_code"`
			FiatCurrencyCode           string      `json:"fiat_currency_code"`
			CryptoCurrencyCode         string      `json:"crypto_currency_code"`
			FiatPricePerCrypto         float64     `json:"fiat_price_per_crypto"`
			FiatAmountRangeMin         int         `json:"fiat_amount_range_min"`
			FiatAmountRangeMax         int         `json:"fiat_amount_range_max"`
			PaymentMethodName          string      `json:"payment_method_name"`
			Active                     bool        `json:"active"`
			PaymentMethodSlug          string      `json:"payment_method_slug"`
			PaymentMethodGroup         string      `json:"payment_method_group"`
			OfferOwnerFeedbackPositive int         `json:"offer_owner_feedback_positive"`
			OfferOwnerFeedbackNegative int         `json:"offer_owner_feedback_negative"`
			LastSeen                   string      `json:"last_seen"`
			LastSeenTimestamp          int         `json:"last_seen_timestamp"`
			RequireVerifiedEmail       bool        `json:"require_verified_email"`
			RequireVerifiedPhone       bool        `json:"require_verified_phone"`
			RequireMinPastTrades       interface{} `json:"require_min_past_trades"`
			RequireVerifiedID          bool        `json:"require_verified_id"`
			PaymentMethodLabel         string      `json:"payment_method_label"`
			OfferTerms                 string      `json:"offer_terms"`
			IsBlocked                  bool        `json:"is_blocked"`
			Tags                       []struct {
				Name        string `json:"name"`
				Slug        string `json:"slug"`
				Description string `json:"description"`
			} `json:"tags"`
			IsFeatured bool `json:"is_featured"`
		} `json:"offers"`
	} `json:"data"`
}
