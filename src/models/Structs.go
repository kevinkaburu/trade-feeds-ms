package models

import (
	"encoding/json"
	"strconv"
)

type SupportedCountries struct {
	//c.country_id,c.name,c.iso_code,fc.fiat_currency_id,fc.currency_code
	CountryID       uint64 `json:"country_id"`
	CountryFlag     string `json:"country_flag"`
	CountryName     string `json:"country_name"`
	CountryCode     string `json:"country_code"`
	FiatCurrencyId  uint64 `json:"fiat_currency_id"`
	FiatCurencyCode string `json:"fiat_currency_code"`
}
type SupportedPaymentMethods struct {
	//payment_method_id,label,payment_type_id
	PaymentMethodId   uint64 `json:"payment_method_id"`
	PaymentMethodName string `json:"payment_method_name"`
	PaymentTypeID     uint64 `json:"payment_type_id"`
}
type SupportedChains struct {
	//chain_name	block_chain_id	chain_code
	ChainName    string `json:"chain_name"`
	BlockChainId string `json:"blockchain_id"`
	ChainCode    string `json:"chain_code"`
}

type SignupPayload struct {
	Email           string    `json:"email"`
	Msisdn          StringInt `json:"msisdn"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
}
type SupportedTokens struct {
	//c.crypto_currency_id,c.name,c.code
	CryptoCurrencyId uint64 `json:"crypto_currency_id"`
	CryptoName       string `json:"crypto_name"`
	CryptoCode       string `json:"crypto_code"`
}

type OfferList struct {
	Count  int            `json:"count"`
	Offers []OfferDbQuery `json:"offers"`
}

type OfferDbQuery struct {
	//crypto_currency_id
	CryptoCurrencyId   int           `json:"crypto_currency_id"`
	ForexExchangeID    int           `json:"forex_exchange_id"`
	ForexExchangeUsd   StringFloat   `json:"usd_exchange_rate"`
	ChainName          string        `json:"chain_name"`
	ChainCode          string        `json:"chain_code"`
	ChainID            string        `json:"chain_id"`
	CountryID          int           `json:"country_id"`
	CountryFlag        string        `json:"country_flag"`
	CountryCode        string        `json:"country_code"`
	CountryName        string        `json:"country_name"`
	OfferID            int           `json:"offer_id"`
	BlockChainId       string        `json:"blockchain_id"`
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
	PaymentMethodId int    `json:"payment_method_id"`
	Tags            string `json:"tags"`
	PaymentMethod   string `json:"payment_method"`
	PaymentType     string `json:"payment_type"`
}

type OfferQuery struct {
	CountryID    StringInt   `json:"country_id"`
	Fiat         string      `json:"fiat_code"`
	FiatAmount   StringFloat `json:"fiat_amount"`
	CryptoAmount StringFloat `json:"crypto_amount"`
}
type StartTradeQuery struct {
	OfferID       StringInt   `json:"offer_id"`
	FiatAmount    StringFloat `json:"fiat_amount"`
	WalletAccount string      `json:"wallet_account"`
	Token         string      `json:"token"`
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

type OfferGetData struct {
	Data struct {
		OfferHash                             string      `json:"offer_hash"`
		ID                                    string      `json:"id"`
		Margin                                int         `json:"margin"`
		Active                                bool        `json:"active"`
		BlockAnonymizerUsers                  bool        `json:"block_anonymizer_users"`
		FiatAmountRangeMin                    int         `json:"fiat_amount_range_min"`
		FiatAmountRangeMax                    int         `json:"fiat_amount_range_max"`
		FeePercentage                         float64     `json:"fee_percentage"`
		CryptoMin                             int         `json:"crypto_min"`
		CryptoMax                             int64       `json:"crypto_max"`
		OfferTerms                            string      `json:"offer_terms"`
		ReleaseTime                           int         `json:"release_time"`
		PaymentMethodLabel                    string      `json:"payment_method_label"`
		PaymentMethodName                     string      `json:"payment_method_name"`
		PaymentMethodSlug                     string      `json:"payment_method_slug"`
		RequireVerifiedEmail                  bool        `json:"require_verified_email"`
		RequireVerifiedPhone                  bool        `json:"require_verified_phone"`
		ShowOnlyTrustedUser                   bool        `json:"show_only_trusted_user"`
		RequireVerifiedID                     bool        `json:"require_verified_id"`
		RequireOfferCurrencyMatchBuyerCountry bool        `json:"require_offer_currency_match_buyer_country"`
		LastSeen                              string      `json:"last_seen"`
		LastSeenTimestamp                     interface{} `json:"last_seen_timestamp"`
		OfferLink                             string      `json:"offer_link"`
		OfferOwnerCountryIso                  interface{} `json:"offer_owner_country_iso"`
		OfferOwnerFeedbackNegative            int         `json:"offer_owner_feedback_negative"`
		OfferOwnerFeedbackPositive            int         `json:"offer_owner_feedback_positive"`
		OfferOwnerProfileLink                 string      `json:"offer_owner_profile_link"`
		OfferOwnerUsername                    string      `json:"offer_owner_username"`
		PaymentWindow                         int         `json:"payment_window"`
		ReleaseTimeMedian                     int         `json:"release_time_median"`
		CurrencyCode                          string      `json:"currency_code"`
		FiatCurrencyCode                      string      `json:"fiat_currency_code"`
		IsBlocked                             bool        `json:"is_blocked"`
		PaymentMethodGroup                    string      `json:"payment_method_group"`
		CryptoCurrency                        string      `json:"crypto_currency"`
		CryptoCurrencyCode                    string      `json:"crypto_currency_code"`
		IsFixedPrice                          bool        `json:"is_fixed_price"`
		BankAccounts                          []struct {
			BankName         string      `json:"bank_name"`
			BankAccountUUID  string      `json:"bank_account_uuid"`
			HolderName       interface{} `json:"holder_name"`
			AccountNumber    interface{} `json:"account_number"`
			FiatCurrencyCode interface{} `json:"fiat_currency_code"`
			IsPersonal       interface{} `json:"is_personal"`
			CountryIso       interface{} `json:"country_iso"`
			SwiftCode        interface{} `json:"swift_code"`
			Iban             interface{} `json:"iban"`
			AdditionalInfo   interface{} `json:"additional_info"`
			RoutingNumber    interface{} `json:"routing_number"`
			Ifsc             interface{} `json:"ifsc"`
			Clabe            interface{} `json:"clabe"`
			BankUUID         interface{} `json:"bank_uuid"`
		} `json:"bank_accounts"`
		FlowType             string      `json:"flow_type"`
		BankReferenceMessage interface{} `json:"bank_reference_message"`
		TradeDetails         string      `json:"trade_details"`
	} `json:"data"`
	Status string `json:"status"`
}

type PaxfulAccount struct {
	Data struct {
		BankName             string      `json:"bank_name"`
		BankAccountUUID      string      `json:"bank_account_uuid"`
		HolderName           string      `json:"holder_name"`
		AccountNumber        string      `json:"account_number"`
		FiatCurrencyCode     string      `json:"fiat_currency_code"`
		IsPersonal           bool        `json:"is_personal"`
		CountryIso           string      `json:"country_iso"`
		SwiftCode            string      `json:"swift_code"`
		Iban                 interface{} `json:"iban"`
		AdditionalInfo       interface{} `json:"additional_info"`
		RoutingNumber        string      `json:"routing_number"`
		Ifsc                 string      `json:"ifsc"`
		Clabe                interface{} `json:"clabe"`
		BankUUID             string      `json:"bank_uuid"`
		InternationalDetails struct {
			Residency interface{} `json:"residency"`
			State     interface{} `json:"state"`
			City      interface{} `json:"city"`
			Zip       interface{} `json:"zip"`
			Address   interface{} `json:"address"`
		} `json:"international_details"`
	} `json:"data"`
	Status    string `json:"status"`
	Timestamp int    `json:"timestamp"`
}

type PaxfulTradeStart struct {
	Status    string `json:"status"`
	Timestamp int    `json:"timestamp"`
	Data      struct {
		Success   bool   `json:"success"`
		TradeHash string `json:"trade_hash"`
	} `json:"data"`
}

type TradeChat struct {
	Status    string `json:"status"`
	Timestamp int    `json:"timestamp"`
	Data      struct {
		Messages []struct {
			ID                string      `json:"id"`
			Timestamp         int         `json:"timestamp"`
			Type              string      `json:"type"`
			TradeHash         string      `json:"trade_hash"`
			IsForModerator    bool        `json:"is_for_moderator"`
			Author            interface{} `json:"author"`
			SecurityAwareness interface{} `json:"security_awareness"`
			Status            int         `json:"status"`
			Text              string      `json:"text"`
			AuthorUUID        interface{} `json:"author_uuid"`
			SentByModerator   bool        `json:"sent_by_moderator"`
		} `json:"messages"`
		Attachments []interface{} `json:"attachments"`
	} `json:"data"`
}

type NewChatResponse struct {
	Status    string `json:"status"`
	Timestamp int    `json:"timestamp"`
	Error     struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type StartTradeResponse struct {
	TradeID       int64          `json:"trade_id"`
	TradeStatus   string         `json:"trade_status"`
	PaymentMethod string         `json:"payment_method"`
	Messages      []ChatMessages `json:"messages"`
}

type ChatMessages struct {
	ID        uint64 `json:"id"`
	Timestamp uint64 `json:"timestamp"`
	Author    string `json:"author"`
	Text      string `json:"text"`
}

type LiveChat struct {
	Token   string    `json:"token"`
	TradeID StringInt `json:"trade_id"`
	Message string    `json:"message"`
}
