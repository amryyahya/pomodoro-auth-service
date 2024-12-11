package main

import (
	"net/http"
	"pomodoro-app/authentication-service/config"
	"pomodoro-app/authentication-service/db"
	"pomodoro-app/authentication-service/models"
	"github.com/gin-gonic/gin"
)

// func getUsers(c *gin.Context) {
//     c.IndentedJSON(http.StatusOK, users)
// }

func register(c *gin.Context) {
    var newUser models.User
    if err := c.BindJSON(&newUser); err != nil {
        return
    }
    config.LoadEnv()
    connStr := config.GetDBConnectionString()
    database := db.Connect(connStr)
    models.InsertUser(database,newUser.Email,[]byte{1,2,3},[]byte{4,5})
	defer database.Close()
    c.IndentedJSON(http.StatusCreated, "hahahaha")
}

func main() {
    router := gin.Default()
    // router.GET("/users", getUsers)
    router.POST("/register", register)


    router.Run("localhost:8000")
}