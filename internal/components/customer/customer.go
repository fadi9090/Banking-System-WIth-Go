package customer

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

func (h *Handler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email is unique (returns true if available)
	isEmailUnique, err := h.repo.UniqueEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isEmailUnique { // Fixed: Check if NOT unique
		c.JSON(http.StatusConflict, gin.H{"error": "Email Does Not exists"}) // Fixed: Use 409 Conflict
		return
	}

	// Check if username is unique
	isUsernameUnique, err := h.repo.UniqueUsername(req.Username) // Fixed: := not =
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isUsernameUnique { // Fixed: Check if NOT unique
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"}) // Fixed: Use 409 Conflict
		return
	}

	customer, err := h.repo.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *Handler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	customer, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *Handler) GetCustomerByNumber(c *gin.Context) {
	number := c.Query("customerNumber")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_number query parameter is required"})
		return
	}

	customer, err := h.repo.GetByCustomerNumber(number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *Handler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var req UpdateCustomerRequest
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if req.Username != "" {
		existing.Username = req.Username
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *Handler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetAllCus(c *gin.Context) {
	customers, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if customers == nil {
		customers = []*Customer{}
	}

	c.JSON(http.StatusOK, customers)
}
