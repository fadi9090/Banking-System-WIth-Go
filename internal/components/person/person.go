package person

import (
	"database/sql"
	"net/http"

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

func (h *Handler) CreatePerson(c *gin.Context) {
	var req CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := h.repo.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &person)
}

func (h *Handler) GetPerson(c *gin.Context) {
	id := c.Param("id")

	person, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if person == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *Handler) GetPersonByEmailOrNumber(c *gin.Context) {
	email := c.Query("email")

	num := c.Query("number")

	if email == "" && num == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email Or Person Number query parameter is required"})
		return
	}

	if email != "" {
		person, err := h.repo.GetByEmail(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if person == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}

		c.JSON(http.StatusOK, person)
	}

	if num != "" {
		person, err := h.repo.GetByNumber(num)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if person == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}

		c.JSON(http.StatusOK, person)
	}
}

func (h *Handler) UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	var req UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.Update(id, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person updated successfully"})
}

func (h *Handler) DeletePerson(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
