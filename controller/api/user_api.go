package api

import (
	"net/http"

	"simple-go-auth/repository/users"
	"simple-go-auth/services/auth"

	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	userId := c.MustGet(auth.ContextUserIdKey).(string)
	if err := users.UpdateUserEmail(userId, user.Email); err != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, err)
	}

	// return Todo struct
	c.JSON(http.StatusOK, "Updated user info")
}

func GetUser(c *gin.Context) {
	userId := c.MustGet(auth.ContextUserIdKey).(string)

	user, err := users.FindUserByID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, err)
	}

	// hide password
	user.Password = "******"

	c.JSON(http.StatusOK, user)
}
