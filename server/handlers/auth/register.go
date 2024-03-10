package auth

import (
	"fmt"
	"time"

	"camera-server/services/database"
	"camera-server/templates"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
    templ.Handler(templates.Register()).ServeHTTP(c.Writer, c.Request)
}

func Create(c *gin.Context) {
    username := c.PostForm("username")
    password := c.PostForm("password")

    db := database.GetDB()

    if err := db.Where("username = ?", username).First(&database.User{}).Error; err == nil {
        values := map[string]string{"username": username, "password": password}

        templ.Handler(templates.RegisterForm(values, map[string]string{"username": "Username is already taken"})).ServeHTTP(c.Writer, c.Request)
        return
    }

    user := database.User{Username: username, Password: password}
    err := db.Create(&user).Error

    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    session := &database.Session{UserID: user.ID}
    db.Create(&session)

    c.SetCookie("session", fmt.Sprint(session.ID), int(time.Hour * 24 * 30), "/", "localhost", false, true)

    c.Header("HX-Redirect", "/")
    c.JSON(200, "success")
}