package address

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

func (h *Handler) CreateAddress(c *gin.Context) {
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := &Address{
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		ZipCode:    req.ZipCode,
		Country:    req.Country,
		CustomerID: req.CustomerID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.repo.Create(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, address)
}

func (h *Handler) GetAddress(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
		return
	}

	address, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if address == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *Handler) GetAddressesByCustomer(c *gin.Context) {
	customerID := c.Param("customerId")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	addresses, err := h.repo.GetByCustomerID(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if addresses == nil {
		addresses = []*Address{}
	}

	c.JSON(http.StatusOK, addresses)
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
		return
	}

	var req UpdateAddressRequest
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if req.Street != "" {
		existing.Street = req.Street
	}
	if req.City != "" {
		existing.City = req.City
	}
	if req.State != "" {
		existing.State = req.State
	}
	if req.ZipCode != "" {
		existing.ZipCode = req.ZipCode
	}
	if req.Country != "" {
		existing.Country = req.Country
	}
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *Handler) DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err.Error() == "address not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
