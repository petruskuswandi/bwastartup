package request

import "github.com/petruskuswandi/bwastartup.git/models"

type GetTransactionCampaignsInput struct {
	ID   string `uri:"id" binding:"required"`
	User models.User
}

type CreateTransactionInput struct {
	Amount     int    `json:"amount" binding:"required"`
	CampaignID string `json:"campaign_id" binding:"required"`
	User       models.User
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
