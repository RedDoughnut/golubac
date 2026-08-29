package main

import (
	//"net/http"
	"context"
	"errors"
	"fmt"
	"os"
	"src/src/auth"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	//"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type SendMessageRequest struct {
	Recipient string
	Message   string
	Token     string
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Server struct {
	DB *pgx.Conn
}
type RegisterRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayname"`
}

func (s *Server) Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user struct {
		id            int64
		password_hash string
		salt          string
	}
	err := s.DB.QueryRow(context.Background(), `SELECT id, password_hash, salt FROM users WHERE username = $1`, input.Username).Scan(&user.id, &user.password_hash, &user.salt)

	if errors.Is(err, pgx.ErrNoRows) {
		err := s.DB.QueryRow(context.Background(), `SELECT id, password_hash, salt FROM users WHERE email = $1`, input.Username).Scan(&user.id, &user.password_hash, &user.salt)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	success, err := auth.VerifyPassword(input.Password, user.password_hash, user.salt)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !success {
		c.JSON(401, gin.H{"error": "Password incorrect"})
		return
	}
	expiresAt := time.Now().AddDate(0, 0, 30) // set the expiry of refresh token to 30 days
	token, err := auth.GenerateRefreshToken()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	hashed_token := auth.HashToken(token)
	_, err = s.DB.Exec(context.Background(), "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3);", user.id, hashed_token, expiresAt)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"refreshtoken": token})
}

// TODO username only alphanumerical + _
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
	var userID int64
	err = s.DB.QueryRow(context.Background(), "INSERT INTO users (username, display_name, email, password_hash, salt) VALUES ($1,$2,$3,$4,$5) RETURNING id;", input.Username, input.DisplayName, input.Email, hashed_pass, salt).Scan(&userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	expiresAt := time.Now().AddDate(0, 0, 30) // set the expiry of refresh token to 30 days
	token, err := auth.GenerateRefreshToken()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	hashed_token := auth.HashToken(token)
	_, err = s.DB.Exec(context.Background(), "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3);", userID, hashed_token, expiresAt)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"refreshtoken": token})
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
		os.Exit(1)
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
	var srv Server
	srv.DB = conn
	defer conn.Close(context.Background())

	router := gin.Default()
	router.POST("/login", srv.Login)
	router.POST("/register", srv.Register)
	router.Run("0.0.0.0:18000")
}
