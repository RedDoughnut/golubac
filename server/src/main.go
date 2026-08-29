package main

import (
	"net/http"
	//"os"
	//"context"
	"fmt"
	//"github.com/jackc/pgx/v5"
	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt/v5"
)

type SendMessageRequest struct {
	Recipient string
	Message   string
	Token     string
}
type LoginRequest struct {
	Username string `json:"username"`
}

func helloworld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}
func Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if input.Username == "" {
		c.JSON(401, gin.H{"error": "No username specified."})
		return
	}

	c.JSON(200, gin.H{"token": "MGnajjacaSkola." + input.Username})
}
func SendMessage(c *gin.Context) {
	var data SendMessageRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if data.Token == "" {
		c.JSON(401, gin.H{"error": "No token specified."})
		return
	}
	c.Status(200)
	fmt.Printf("Got a successful request! token: %s\nrecepient: %s\nmessage: %s\n", data.Token, data.Recipient, data.Message)
}
func main() {
	// ctx := context.Background()
	// targetdb := "golubac"
	// username := "admin"
	// password := "admin"
	// host := "localhost"
	// port := "18000"

	router := gin.Default()
	router.GET("/helloworld", helloworld)
	router.POST("/login", Login)
	router.POST("/send-message", SendMessage)
	router.Run("0.0.0.0:18000")
}
