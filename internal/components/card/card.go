package card

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateCard(c *gin.Context) {
	accountID := c.Param("accountId")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.repo.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, card)
}

func (h *Handler) GetCard(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card ID is required"})
		return
	}

	card, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if card == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *Handler) GetCardsByAccount(c *gin.Context) {
	accountID := c.Param("accountId")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	cards, err := h.repo.GetByAccountID(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cards == nil {
		cards = []*Card{}
	}

	c.JSON(http.StatusOK, cards)
}

func (h *Handler) UpdateCard(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card ID is required"})
		return
	}

	var req UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	if req.Brand != "" {
		existing.Brand = req.Brand
	}
	if req.IsBlocked != nil {
		existing.IsBlocked = *req.IsBlocked
	}
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *Handler) BlockCard(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card ID is required"})
		return
	}

	var req BlockCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Block(id, true); err != nil {
		if err.Error() == "card not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card blocked successfully"})
}
