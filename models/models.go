package models

import (
	"time"
)

// Transaction Processing details
type ClientInfo struct {
	ServiceCode           string              `json:"ServiceCode"  binding:"required"`
	AppCode               string              `json:"AppCode"  binding:"required"`
	Settings              map[string]string   `json:"Settings"  binding:"required"`
	SettingsDb            string              `json:"-"`
	Secret                string              `json:"Secret"  binding:"required"`
	Multinational         bool                `json:"Multinational"  binding:"required"`
	StatusDesciption      string              `json:"StatusDesciption"  binding:"required"`
	StatusCode            string              `json:"StatusCode"  binding:"required"`
	AppStatusCode         string              `json:"AppStatusCode"  binding:"required"`
	AppStatusCodeDesc     string              `json:"AppStatusCodeDesc"  binding:"required"`
	ServiceStatusCodeDesc string              `json:"ServiceStatusCodeDesc"  binding:"required"`
	ServiceStatusCode     string              `json:"ServiceStatusCode"  binding:"required"`
	CurrencyCode          string              `json:"CurrencyCode"  binding:"required"`
	Headers               map[string][]string `json:"Headers"  binding:"required"`
	HeadersDb             string              `json:"-"`
	Statics               map[string]string   `json:"Statics"  binding:"required"`
	StaticsDb             string              `json:"-"`
	HTTPMethod            string              `json:"HTTPMethod"  binding:"required"`
	ServiceURL            string              `json:"ServiceURL"  binding:"required"`
	ServiceCountry        uint                `json:"ServiceCountry"  binding:"required"`
	AppCountry            uint                `json:"AppCountry"  binding:"required"`
}

// Incomming payload from Payments Queue
type Request struct {
	Transaction Transaction `json:"transaction" binding:"required"`
	ClientInfo  ClientInfo  `json:"settings" binding:"-"`
}

// Entire transaction details
type Transaction struct {
	MSISDN          string                 `json:"MSISDN" binding:"required"`
	AccountNumber   string                 `json:"accountNumber" binding:"required"`
	ExternalCode    string                 `json:"externalCode" binding:"required"`
	Amount          string                 `json:"amount" binding:"required"`
	PaymentDate     time.Time              `json:"currentDate" binding:"required"`
	Narration       string                 `json:"narration" binding:"required"`
	ISOCurrencyCode string                 `json:"ISOCurrencyCode" binding:"required"`
	CustomerName    string                 `json:"customerName" binding:"required"`
	PaymentMode     string                 `json:"paymentMode" binding:"required"`
	Callback        string                 `json:"callback" binding:"required"`
	Code            string                 `json:"Code" binding:"required"`
	Metadata        map[string]interface{} `json:"metadata" binding:"-"`
	AppKey          string                 `json:"appKey" binding:"-"`
}

// Response to be forwarded to payments.callbackXchange for sync
// Status code only one changing depending on async or sync
type Response struct {
	ExternalCode      string    `json:"transactionId" binding:"required"`
	Code              string    `json:"trackingId" binding:"required"`
	RecievedDate      time.Time `json:"recievedDate" binding:"required"`
	StatusCode        string    `json:"statusCode" binding:"required"`
	StatusDescription string    `json:"statusDescription" binding:"required"`
}

// JS Error defination and details
type Error struct {
	Message    string `json:"message" binding:"-"`
	Location   string `json:"location" binding:"-"`
	StackTrace string `json:"stacktrace" binding:"-"`
}

// Returned request from JS Script
type RequestBuilt struct {
	Payload string              `json:"payload" binding:"payload"`
	Headers map[string][]string `json:"headers" binding:"required"`
	Error   string              `json:"error" binding:"-"`
}
