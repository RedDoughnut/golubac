package main

import (
	"net/http"
	//"github.com/jackc/pgx/v5"
	"github.com/gin-gonic/gin"
)

func helloworld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}
func main() {
	router := gin.Default()
	router.GET("/helloworld", helloworld)
	router.Run("localhost:18000")
}
