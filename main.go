package main

import (
    "database/sql"
	"fmt"
	"net/http"
	"pomodoro-app/authentication-service/config"
	"pomodoro-app/authentication-service/db"
	"pomodoro-app/authentication-service/models"
	"pomodoro-app/authentication-service/utils"
	"github.com/gin-gonic/gin"
)

var dbPool *sql.DB

func init(){
    config.LoadEnv()
    connStr := config.GetDBConnectionString()
    dbPool = db.Connect(connStr)
}


func login(c *gin.Context) {
    var user models.User
    if err := c.BindJSON(&user); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
        return
    }
    database := dbPool
    isUserExisted, err := models.CheckIfEmailExist(database, user.Email)
    if err != nil {
        fmt.Print(err)
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
        return
    }
    if !isUserExisted {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error":"User not Found"})
        return
    }
    hashedPassword, salt, err := models.GetUserCredByEmail(database, user.Email)
    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
        return
    }
    isValid := utils.VerifyPassword(user.Password,hashedPassword,salt)
    if isValid {
        accessSecret,refreshSecret := config.GetJWTSecret()
        user_id, err := models.GetUserIDByEmail(database, user.Email)
        if err != nil {
            fmt.Println(err)
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
            return
        }
        accessToken, refreshToken, err := utils.GenerateTokens(user_id,accessSecret,refreshSecret)
        if err!=nil {
            fmt.Println(err)
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
            return
        }
        c.IndentedJSON(http.StatusAccepted, gin.H{"access-token":accessToken,"refresh-token":refreshToken})
    } else {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password"})
    }
}

func register(c *gin.Context) {
    var newUser models.User
    if err := c.BindJSON(&newUser); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
        return
    }
    database := dbPool
    isUserExisted, _ := models.CheckIfEmailExist(database, newUser.Email)
	if isUserExisted {
		c.IndentedJSON(http.StatusConflict, gin.H{"error":"Email Has Been Registered"})
        return
	}

    hashedPassword, salt, _ := utils.HashPassword(newUser.Password)
    err := models.InsertUser(database,newUser.Email,hashedPassword,salt)
    if err!=nil {
        fmt.Print(err)
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
        return
    } else {
        accessSecret,refreshSecret := config.GetJWTSecret()
        user_id, err := models.GetUserIDByEmail(database, newUser.Email)
        if err != nil {
            fmt.Println(err)
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
            return
        }
        accessToken, refreshToken, err := utils.GenerateTokens(user_id,accessSecret,refreshSecret)
        if err!=nil {
            fmt.Println(err)
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
            return
        }
        c.IndentedJSON(http.StatusCreated, gin.H{"access-token":accessToken,"refresh-token":refreshToken})
    }
    
}

func refresh(c *gin.Context) {
    var requestData map[string]interface{}

    // Bind the JSON request body to a map
    if err := c.BindJSON(&requestData); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
        return
    }
    refreshToken := requestData["refresh-token"].(string)
    accessSecret,refreshSecret := config.GetJWTSecret()
    newAccessToken, err := utils.RefreshAccessToken(refreshToken, accessSecret, refreshSecret)
    if err != nil{
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
        return
    } else {
        c.IndentedJSON(http.StatusCreated, gin.H{"access-token":newAccessToken})
    }

}

func validate(c *gin.Context) {
    var requestData map[string]interface{}

    if err := c.BindJSON(&requestData); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
        return
    }
    database := dbPool

    refreshToken := requestData["refresh-token"].(string)
    tokenInvalid, err:= utils.IsTokenBlacklisted(database, refreshToken)
    if err != nil {
        fmt.Println(err)
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
        return
    }
    if tokenInvalid {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Token is Invalid"})
        return
    }
    c.IndentedJSON(http.StatusAccepted, gin.H{"message": "Valid"})
}

func logout(c *gin.Context) {
    var requestData map[string]interface{}

    // Bind the JSON request body to a map
    if err := c.BindJSON(&requestData); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
        return
    }
    refreshToken := requestData["refresh-token"].(string)
    _,refreshSecret := config.GetJWTSecret()
    claims,err := utils.ValidateToken(refreshToken, []byte(refreshSecret))

    if err != nil {
        fmt.Println(err)
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
        return
    }
    database := dbPool
    err = utils.BlacklistToken(database, claims["user_id"].(string), refreshToken, claims["exp"].(float64))
    if err != nil {
        fmt.Println(err)
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
        return
    }
    c.IndentedJSON(http.StatusAccepted, gin.H{"message":"logout success"})
}

func main() {
    router := gin.Default()

    router.POST("/api/v1/login", login)
    router.POST("/api/v1/register", register)
    router.POST("/api/v1/login/refresh", refresh)
    router.POST("/api/v1/login/validate", validate)
    router.POST("/api/v1/logout", logout)
    router.Run("localhost:8000")
}