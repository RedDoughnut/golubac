package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func login(c *gin.Context){
	c.IndentedJSON(http.StatusOK, "Hello World!")
}
func main() {
	
}
