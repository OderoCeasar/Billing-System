package mpesa

import "time"


type TokenResponse struct {
	AccessToken  string	 `json:"access_token"`
	ExpiresIn  	 string  `json:"expires_in"`
}

// STK push payload request
type STKPushRequest struct {
	BusinessShortCode   string	`json:"BusinessShortCode"`
	Password			string	`json:"Password"`
	Timestamp			string	`json:"Timestamp"`
	TransactionType		string	`json:"TransactionType"`
	Amount				string	`json:"Amount"`
	PartyA				string	`json:"PartyA"`
	PartyB				string	`json:"PartyB"`
	PhoneNumber			string	`json:"PhoneNumber"`
	CallbackURL			string	`json:"CallbackURL"`
	AccountReference 	string  `json:"AccountReference"`
	TransactionDesc		string	`json:"TransactionDesc"`
}

// STK push response
type STKPushResponse struct {
	MerchantRequestID	string	`json:"MerchantRequestID"`
	CheckoutRequestID 	string	`json:"CheckoutRequestID"`
	ResponseCode		string 	`json:"ResponseCode"`
	ResponseDescription	string 	`json:"ResponseDescription"`
	CustomerMessage		string  `json:"CustomerMessage"`
}

type CallbackMetadataItem struct {
	Name 	string	`json:"Name"`
	Value 	interface{} `json:"Value"`
}

type CallbackMetadata struct {
	Item []CallbackMetadataItem	`json:"Item"`
}

type StkCallback struct {
	MerchantRequestID	string			  `json:"MerchantRequestID"`
	CheckoutRequestID 	string			  `json:"CheckoutRequestID"`
	ResultCode			int 			  `json:"ResultCode"`
	ResultDesc			string 			  `json:"ResultDesc"`
	CallbackMetadata	CallbackMetadata  `json:"CallbackMetadata"`
}

type CallbackBody struct {
	StkCallback StkCallback	`json:"stkcallback"`
}

type MpesaCallback struct {
	Body CallbackBody `json:"Body"`
}

// Parsed Callback
type ParsedCallback struct {
	CheckoutRequestID	string
	ResultCode			int
	ResultDesc			string
	Amount				float64
	MpesaReceiptNumber	string
	TransactionDate		time.Time
	PhoneNumber			string
}