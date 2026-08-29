package main

import (
	//"net/http"
	"context"
	"fmt"
	"os"

	"src/src/auth"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	//"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type SendMessageRequest struct {
	Recipient string
	Message   string
	Token     string
}
type LoginRequest struct {
	Username string `json:"username"`
}
type Server struct {
	DB pgx.Conn
}
type RegisterRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayname"`
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
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

func (s *Server) Register(c *gin.Context) {
	var input RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if input.Email == "" || input.DisplayName == "" || input.Username == "" || input.Password == "" {
		c.JSON(400, gin.H{"error": "Missing fields."})
		return
	}
	salt, err := auth.GenerateSalt()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	hashed_pass := auth.HashPassword(input.Password, salt)
	_,err = s.DB.Exec(context.Background(), "INSERT INTO users (username, display_name, email, hashed_password, salt) VALUES ($1,$2,$3,$4,$5);", input.Username, input.DisplayName, input.Email, hashed_pass, salt)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"refresh"})
}
func SendTextMessage(c *gin.Context) {
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
func printError(format string, a ...interface{}) {
	fmt.Printf("\033[1;31m"+format+"\033[0m\n", a...)
}

// Loads enviromental variables
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		printError("Error loading .env file")
	}
}
func main() {
	loadEnv()
	PostgresURL := os.Getenv("POSTGRES_URL")
	conn, err := pgx.Connect(context.Background(), PostgresURL)
	if err != nil {
		printError("Error connecting to DB: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	router := gin.Default()
	router.POST("/login", Login)
	router.POST("/send-text-message", SendMessage)
	router.Run("0.0.0.0:18000")
}
