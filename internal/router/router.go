package router

import (
	"strings"
	"task-project/internal/components/account"
	"task-project/internal/components/address"
	"task-project/internal/components/card"
	"task-project/internal/components/customer"
	"task-project/internal/components/person"
	transaction "task-project/internal/components/transactions"
	"task-project/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	personHandler *person.Handler,
	customerHandler *customer.Handler,
	addressHandler *address.Handler,
	accountHandler *account.Handler,
	cardHandler *card.Handler,
	transactionHandler *transaction.Handler,
) *gin.Engine {

	r := gin.Default()

	// Health check (no authentication)
	r.GET("/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API group with authentication
	api := r.Group("/v1")
	api.Use(authMiddleware(cfg.APIKey))
	{
		// ===== PERSON ROUTES =====
		api.POST("/person", personHandler.CreatePerson)
		api.GET("/person/:id", personHandler.GetPerson)
		api.GET("/person", personHandler.GetPersonByEmailOrNumber)
		api.PATCH("/person/:id", personHandler.UpdatePerson)
		api.DELETE("/person/:id", personHandler.DeletePerson)

		// ===== CUSTOMER ROUTES =====
		api.POST("/customer", customerHandler.CreateCustomer)
		api.GET("/customer/:id", customerHandler.GetCustomer)
		api.GET("/customer", customerHandler.GetCustomerByNumber)
		api.GET("/customers", customerHandler.GetAllCus)
		api.PATCH("/customer/:id", customerHandler.UpdateCustomer)
		api.DELETE("/customer/:id", customerHandler.DeleteCustomer)

		// ===== ADDRESS ROUTES (nested under customers) =====
		api.POST("/customer/:id/addresses", addressHandler.CreateAddress)
		api.GET("/customer/:id/addresses", addressHandler.GetAddressesByCustomer)
		api.PATCH("/address/:address_id", addressHandler.UpdateAddress)
		api.DELETE("/address/:address_id", addressHandler.DeleteAddress)

		// ===== ACCOUNT ROUTES =====
		api.POST("/customer/:id/accounts", accountHandler.CreateAccount)
		api.GET("/account/:id", accountHandler.GetAccount)
		api.GET("/customer/:id/accounts", accountHandler.GetAccountsByCustomer)
		api.PATCH("/account/:id", accountHandler.UpdateAccount)
		api.POST("/account/:id/close", accountHandler.CloseAccount)

		// ===== CARD ROUTES (nested under accounts) =====
		api.POST("/account/:id/cards", cardHandler.CreateCard)
		api.GET("/account/:id/cards", cardHandler.GetCardsByAccount)
		api.GET("/card/:card_id", cardHandler.GetCard)
		api.PATCH("/card/:card_id", cardHandler.UpdateCard)
		api.POST("/card/:card_id/block", cardHandler.BlockCard)

		// ===== TRANSACTION ROUTES (nested under accounts) =====
		api.POST("/account/:id/transactions", transactionHandler.CreateTransaction)
		api.GET("/account/:id/transactions", transactionHandler.GetTranByAccId)
		api.GET("/transaction/:tx_id", transactionHandler.GetTranById)
	}

	return r
}

func authMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token == "" {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"message": "Token is required",
			})
			c.Abort()
			return
		}

		if token != apiKey {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
