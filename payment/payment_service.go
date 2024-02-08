package payment

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/petruskuswandi/bwastartup.git/models"
	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user models.User) (string, error)
}

func NewService() *service {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user models.User) (string, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")

	if serverKey == "" || clientKey == "" {
		log.Println("Midtrans API keys are not set. Please check your .env file.")
		return "", fmt.Errorf("midtrans API keys are not set")
	}

	midclient := midtrans.NewClient()
	midclient.ServerKey = serverKey
	midclient.ClientKey = clientKey
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.ID,
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}
