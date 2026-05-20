package utility

import (
	"github.com/gin-gonic/gin"
)

// RFC7807Error represents the standard error format required by the assignment
type RFC7807Error struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

// SendSuccess sends a successful JSON response
func SendSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// SendError sends an RFC7807 formatted error response
func SendError(c *gin.Context, status int, title string, detail string) {
	err := RFC7807Error{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	c.JSON(status, err)
}
