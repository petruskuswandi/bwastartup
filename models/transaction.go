package models

import (
	"time"

	"github.com/leekchan/accounting"
	"gorm.io/gorm"
)

type Transaction struct {
	ID         string
	CampaignID string
	UserID     string
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       User
	Campaign   Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func (c Transaction) AmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(c.Amount)
}
