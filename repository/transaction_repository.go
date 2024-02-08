package repository

import (
	"github.com/petruskuswandi/bwastartup.git/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetByCampaignID(campaignID string) ([]models.Transaction, error)
	GetByUserID(userID string) ([]models.Transaction, error)
	GetByID(ID string) (models.Transaction, error)
	Save(transaction models.Transaction) (models.Transaction, error)
	GetCurrentCounter() (int64, error)
	Update(transaction models.Transaction) (models.Transaction, error)
}
type transactionRepository struct {
	db *gorm.DB
}

func NewRepositoryTransaction(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetByCampaignID(campaignID string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetByUserID(userID string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetByID(ID string) (models.Transaction, error) {
	var transaction models.Transaction

	err := r.db.Where("id = ?", ID).Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) Save(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) Update(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Save(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) GetCurrentCounter() (int64, error) {
	var count int64

	result := r.db.Table("transactions").Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}
