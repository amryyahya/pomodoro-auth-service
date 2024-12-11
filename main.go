package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "pomodoro-app/authentication-service/config"
	"pomodoro-app/authentication-service/db"
)

type user struct {
    Email           string  `json:"email"`
    HashedPassword  string  `json:"hashedPassword"`
}

var users = []user{
    {Email: "Blue Train", HashedPassword: "password"},
    {Email: "Jeru", HashedPassword: "password"},
    {Email: "Sarah Vaughan and Clifford Brown", HashedPassword: "password"},
}

func getUsers(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, users)
}

func register(c *gin.Context) {
    var newUser user
    if err := c.BindJSON(&newUser); err != nil {
        return
    }
    users = append(users, newUser)
    c.IndentedJSON(http.StatusCreated, newUser)
}

func main() {
    router := gin.Default()
    router.GET("/users", getUsers)
    router.POST("/register", register)
    config.LoadEnv()
    connStr := config.GetDBConnectionString()
    database := db.Connect(connStr)
	defer database.Close()

    router.Run("localhost:8080")
}