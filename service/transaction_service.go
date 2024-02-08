package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/petruskuswandi/bwastartup.git/models"
	"github.com/petruskuswandi/bwastartup.git/payment"
	"github.com/petruskuswandi/bwastartup.git/repository"
	"github.com/petruskuswandi/bwastartup.git/request"
)

type serviceTransaction struct {
	repository         repository.TransactionRepository
	campaignRepository repository.CampaignRepository
	paymentService     payment.Service
	mutex              sync.Mutex
}
type ServiceTransaction interface {
	GetTransactionByCampaignID(input request.GetTransactionCampaignsInput) ([]models.Transaction, error)
	GetTransactionByUserID(userID string) ([]models.Transaction, error)
	CreateTransaction(input request.CreateTransactionInput) (models.Transaction, error)
	ProcessPayment(input request.TransactionNotificationInput) error
	GetAllTransactions() ([]models.Transaction, error)
}

func NewServiceTransaction(repo repository.TransactionRepository, campaignRepository repository.CampaignRepository, paymentService payment.Service, mutex sync.Mutex) *serviceTransaction {
	return &serviceTransaction{repo, campaignRepository, paymentService, mutex}
}

func (s *serviceTransaction) GetTransactionByCampaignID(input request.GetTransactionCampaignsInput) ([]models.Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []models.Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []models.Transaction{}, errors.New("not an owner of the campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *serviceTransaction) GetTransactionByUserID(userID string) ([]models.Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *serviceTransaction) CreateTransaction(input request.CreateTransactionInput) (models.Transaction, error) {
	transaction := models.Transaction{}
	transactionID, err := uuid.NewRandom()
	if err != nil {
		return transaction, err
	}
	transaction.ID = transactionID.String()
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"

	code, err := s.generateTransactionCode()
	if err != nil {
		return models.Transaction{}, err
	}

	transaction.Code = code

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, nil
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, nil
	}

	return newTransaction, nil
}

func (s *serviceTransaction) ProcessPayment(input request.TransactionNotificationInput) error {
	transaction_id := input.OrderID

	transaction, err := s.repository.GetByID(transaction_id)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *serviceTransaction) generateTransactionCode() (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentCounter, err := s.getCurrentCounterFromDatabase()
	if err != nil {
		return "", err
	}

	nextCounter := currentCounter + 1

	err = s.saveCounterToDatabase(nextCounter)
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf("ORDER-%03d", nextCounter)

	return code, nil
}

func (s *serviceTransaction) getCurrentCounterFromDatabase() (int64, error) {
	counter, err := s.repository.GetCurrentCounter()

	if err != nil {
		return 0, err
	}

	return counter, nil
}

func (s *serviceTransaction) saveCounterToDatabase(counter int64) error {
	return nil
}

func (s *serviceTransaction) GetAllTransactions() ([]models.Transaction, error) {
	transactions, err := s.repository.FindAll()
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}