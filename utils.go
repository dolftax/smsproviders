package smsproviders

import (
    "time"
    "net/http"
)


type DeliveryStatus string

const (
    DeliveryStatusPending DeliveryStatus = "PENDING"
    DeliveryStatusAccepted DeliveryStatus = "ACCEPTED"
    DeliveryStatusFailed DeliveryStatus = "FAILED"
    DeliveryStatusDelivered DeliveryStatus = "DELIVERED"
)

type SMSStatusReport struct {
    Provider        string
    MessageId       string
    CustomId        string
    ProviderStatus  string
    DeliveryStatus  DeliveryStatus
    Description     string
    DeliveryTime    *time.Time
}

type SMSProvider interface {
	SendSMS(string, string, string) (*SMSStatusReport, error)
	GetStatusReport(string) (*SMSStatusReport, error)
	CallbackToReport(*http.Request) (*SMSStatusReport, error)
}
