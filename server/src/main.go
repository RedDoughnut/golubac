package main

import (
	//"net/http"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"src/src/auth"
	"src/src/ratelimiter"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"

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
	DB              *pgx.Conn
	JWTSecret       string
	WSHub           *Hub
	Upgrader        *websocket.Upgrader
	RegisterLimiter *ratelimiter.RateLimiter
	LoginLimiter    *ratelimiter.RateLimiter
}
type RegisterRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayname"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refreshtoken"`
}

// Hub -> WritePump
type OutgoingMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Websocket -> ReadPump
type IncomingMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
	To   int64  `json:"to"`
}

// ReadPump -> Hub
type HubMessage struct {
	Text string
	To   int64
}
type Client struct {
	UserID int64
	Conn   *websocket.Conn

	Send chan OutgoingMessage
}

type Hub struct {
	clients map[int64]map[*Client]bool

	register   chan *Client
	unregister chan *Client
	send       chan *HubMessage
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
		case client := <-h.unregister:
			delete(h.clients[client.UserID], client)
			if len(h.clients[client.UserID]) == 0 {
				delete(h.clients, client.UserID)
			}
			close(client.Send)
		case broadcastMessage := <-h.send:
			for client, _ := range h.clients[broadcastMessage.To] {
				client.Send <- OutgoingMessage{
					Text: broadcastMessage.Text,
					Type: "message",
				}
			}
		}
	}
}

func (c *Client) WritePump() {
	for message := range c.Send {
		if err := c.Conn.WriteJSON(message); err != nil {
			return
		}
	}
}

// TODO: Read username from websocket instead of userID
// TODO: Save messages to Postgres
func (c *Client) ReadPump(h *Hub) {
	defer func() {
		c.Conn.Close()
		h.unregister <- c
	}()

	for {
		var message IncomingMessage
		if err := c.Conn.ReadJSON(&message); err != nil {
			return
		}
		h.send <- &HubMessage{
			Text: message.Text,
			To:   message.To,
		}
	}
}
func (s *Server) Login(c *gin.Context) {
	if !s.LoginLimiter.Allow(c.ClientIP()) {
		c.JSON(429, gin.H{"error": "Too many login attempts!"})
		return
	}
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
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(401, gin.H{"error": "Invalid username/email or password"})
			return
		} else if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	success, err := auth.VerifyPassword(input.Password, user.password_hash, user.salt)
	if err != nil || !success {
		c.JSON(401, gin.H{"error": "Invalid username/email or password"})
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
	if !s.RegisterLimiter.Allow(c.ClientIP()) {
		c.JSON(429, gin.H{"error": "Too many registrations!"})
		return
	}
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
	password_hash, err := auth.HashPassword(input.Password, salt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var userID int64
	err = s.DB.QueryRow(context.Background(), "INSERT INTO users (username, display_name, email, password_hash, salt) VALUES ($1,$2,$3,$4,$5) RETURNING id;", input.Username, input.DisplayName, input.Email, password_hash, salt).Scan(&userID)
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

// TODO dodati limiter
func (s *Server) Refresh(c *gin.Context) {
	var input RefreshRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	hashed_token := auth.HashToken(input.RefreshToken)
	var user struct {
		id         int64
		expires_at time.Time
		revoked_at sql.NullTime
	}
	err := s.DB.QueryRow(context.Background(), "SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token_hash = $1", hashed_token).Scan(&user.id, &user.expires_at, &user.revoked_at)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// check if token is expired or revoked
	if user.expires_at.Before(time.Now()) || user.revoked_at.Valid {
		c.JSON(401, gin.H{"error": "Refresh token expired"})
		return
	}
	session_token, err := auth.GenerateJWT(user.id, s.JWTSecret)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"sessiontoken": session_token})
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

func (srv *Server) rateLimiterCleanup() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		srv.LoginLimiter.Cleanup()
		srv.RegisterLimiter.Cleanup()
	}
}
func (s *Server) HandleWebSocket(c *gin.Context) {
	token := c.GetHeader("Authorization")

	valid, userID, errType := auth.VerifyJWT(token, s.JWTSecret)
	if !valid {
		// JWT expired
		if errType == 1 {
			c.JSON(401, gin.H{"error": "Session token expired"})
		} else {
			c.JSON(400, gin.H{"error": "Invalid session token."})
		}
		return
	}
	conn, err := s.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan OutgoingMessage, 256),
	}

	s.WSHub.register <- client

	go client.ReadPump(s.WSHub)
	go client.WritePump()
}

func main() {
	loadEnv()
	PostgresURL := os.Getenv("POSTGRES_URL")
	JWTSecret := os.Getenv("JWT_SECRET")
	Port := os.Getenv("PORT")
	conn, err := pgx.Connect(context.Background(), PostgresURL)
	if err != nil {
		printError("Error connecting to DB: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var srv Server
	srv.DB = conn
	srv.JWTSecret = JWTSecret
	srv.RegisterLimiter = ratelimiter.New(10, 24*time.Hour) // set rate limit for register to 10/day
	srv.LoginLimiter = ratelimiter.New(5, time.Hour)        // set rate limit for login 5/h
	go srv.rateLimiterCleanup()
	// WEB SOCKET
	WSHub := Hub{
		clients:    make(map[int64]map[*Client]bool),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		send:       make(chan *HubMessage, 256),
	}
	srv.WSHub = &WSHub
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	srv.Upgrader = &upgrader

	go srv.WSHub.Run()

	router := gin.Default()
	router.POST("/login", srv.Login)
	router.POST("/register", srv.Register)
	router.POST("/refresh-session-token", srv.Refresh)
	router.GET("/ws", srv.HandleWebSocket)
	router.Run("localhost:" + Port)
}
