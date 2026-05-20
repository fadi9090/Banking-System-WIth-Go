package transaction

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateTransaction(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})
}

func (h *Handler) GetTranById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	tran, err := h.repo.GetTranById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tran == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, tran)
}

func (h *Handler) GetTranByAccId(c *gin.Context) {
	id := c.Param("accountId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	trans, err := h.repo.GetTranByAccId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trans == nil {
		trans = []*Transaction{}
	}

	c.JSON(http.StatusOK, trans)
}

func (h *Handler) GetTranByAmount(c *gin.Context) {
	amount := c.Param("amount")
	if amount == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount is required"})
		return
	}

	amountNumber, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	trans, err := h.repo.GetTranByAmount(amountNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trans == nil {
		trans = []*Transaction{}
	}

	c.JSON(http.StatusOK, trans)
}

func (h *Handler) GetTranByCardId(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card ID is required"})
		return
	}

	trans, err := h.repo.GetTranByCardId(cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trans == nil {
		trans = []*Transaction{}
	}

	c.JSON(http.StatusOK, trans)
}

func (h *Handler) GetTranByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	trans, err := h.repo.GetTranByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trans == nil {
		trans = []*Transaction{}
	}

	c.JSON(http.StatusOK, trans)
}

func (h *Handler) GetTranByDirection(c *gin.Context) {
	dir := c.Param("direction")
	if dir == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Direction is required"})
		return
	}

	trans, err := h.repo.GetTranByDirection(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trans == nil {
		trans = []*Transaction{}
	}

	c.JSON(http.StatusOK, trans)
}

func (h *Handler) GetTranByDate(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date is required"})
		return
	}

	dateFormat, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	trans, err := h.repo.GetTranByDate(dateFormat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trans == nil {
		trans = []*Transaction{}
	}

	c.JSON(http.StatusOK, trans)
}

func (h *Handler) GetAllTransactions(c *gin.Context) {
	trans, err := h.repo.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trans == nil {
		trans = []*Transaction{}
	}

	c.JSON(http.StatusOK, trans)
}
