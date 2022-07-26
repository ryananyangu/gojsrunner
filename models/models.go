package models

import (
	"time"
)

type Request struct {
	Headers       map[string][]string `json:"headers" binding:"required"`
	Payload       Payload             `json:"payload" binding:"required"`
	Constants     map[string]string   `json:"constants" binding:"required"`
	ClientConfigs map[string]string   `json:"AppSettings" binding:"required"`
}

type Payload struct {
	MSISDN          string                 `json:"MSISDN" binding:"required"`
	AccountNumber   string                 `json:"accountNumber" binding:"required"`
	TransactionId   string                 `json:"transactionId" binding:"required"`
	Amount          string                 `json:"amount" binding:"required"`
	CurrentDate     time.Time              `json:"currentDate" binding:"required"`
	Narration       string                 `json:"narration" binding:"required"`
	ClientCode      string                 `json:"clientCode" binding:"required"`
	ISOCurrencyCode string                 `json:"ISOCurrencyCode" binding:"required"`
	CustomerName    string                 `json:"customerName" binding:"required"`
	PaymentMode     string                 `json:"paymentMode" binding:"required"`
	Callback        string                 `json:"callback" binding:"required"`
	Metadata        map[string]interface{} `json:"metadata" binding:"-"`
}

type Response struct {
	TransactionId     string                 `json:"transactionId" binding:"required"`
	TrackingId        string                 `json:"trackingId" binding:"required"`
	RecievedDate      time.Time              `json:"recievedDate" binding:"required"`
	StatusCode        string                 `json:"statusCode" binding:"required"`
	StatusDescription string                 `json:"statusDescription" binding:"required"`
	Metadata          map[string]interface{} `json:"metadata" binding:"-"`
}

type Error struct {
	Message    string `json:"message" binding:"-"`
	Location   string `json:"location" binding:"-"`
	StackTrace string `json:"stacktrace" binding:"-"`
}

type RequestBuilt struct {
	Payload string            `json:"payload" binding:"payload"`
	Headers map[string]string `json:"headers" binding:"required"`
	Error   Error             `json:"error" binding:"-"`
}

type ResponseBuilt struct {
	Response map[string]interface{} `json:"response" binding:"required"`
}
