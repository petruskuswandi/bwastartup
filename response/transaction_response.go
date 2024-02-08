package response

import (
	"time"

	"github.com/petruskuswandi/bwastartup.git/models"
)

type CampaignTransactionResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func ResponseCampaignTransaction(transaction models.Transaction) CampaignTransactionResponse {
	formatter := CampaignTransactionResponse{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}

func ResponseCampaignTransactions(transactions []models.Transaction) []CampaignTransactionResponse {
	if len(transactions) == 0 {
		return []CampaignTransactionResponse{}
	}

	var transactionFormatter []CampaignTransactionResponse

	for _, transaction := range transactions {
		formatter := ResponseCampaignTransaction(transaction)
		transactionFormatter = append(transactionFormatter, formatter)
	}

	return transactionFormatter
}

type UserTransactionResponse struct {
	ID        string            `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction models.Transaction) UserTransactionResponse {
	formatter := UserTransactionResponse{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	CampaignFormatter := CampaignFormatter{}
	CampaignFormatter.Name = transaction.Campaign.Name
	CampaignFormatter.ImageURL = ""
	if len(transaction.Campaign.CampaignImages) > 0 {
		CampaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	formatter.Campaign = CampaignFormatter

	return formatter
}

func ResponseUserTransactions(transactions []models.Transaction) []UserTransactionResponse {
	if len(transactions) == 0 {
		return []UserTransactionResponse{}
	}

	var transactionsFormatter []UserTransactionResponse

	for _, transaction := range transactions {
		formatter := FormatUserTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}

type TransactionFormatter struct {
	ID         string `json:"id"`
	CampaignID string `json:"campaign_id"`
	UserID     string `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

func FormatTransaction(transaction models.Transaction) TransactionFormatter {
	formatter := TransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.CampaignID = transaction.CampaignID
	formatter.UserID = transaction.UserID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.Code = transaction.Code
	formatter.PaymentURL = transaction.PaymentURL
	return formatter
}
