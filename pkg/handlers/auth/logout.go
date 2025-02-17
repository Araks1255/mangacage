package auth

import "github.com/gin-gonic/gin"

func (h handler) Logout(c *gin.Context) {
	c.SetCookie("mangabrad_token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"error":"Выход из аккаунта выполнен успешно"})
}