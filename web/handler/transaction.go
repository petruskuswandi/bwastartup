package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petruskuswandi/bwastartup.git/service"
)

type transactionHandler struct {
	transactionService service.ServiceTransaction
}

func NewTransactionHandler(transactionService service.ServiceTransaction) *transactionHandler {
	return &transactionHandler{transactionService}
}

func (h *transactionHandler) Index(c *gin.Context) {
	transactions, err := h.transactionService.GetAllTransactions()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "transaction_index.html", gin.H{"transactions": transactions})
}
