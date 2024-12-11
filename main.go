package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
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

    router.Run("localhost:8080")
}