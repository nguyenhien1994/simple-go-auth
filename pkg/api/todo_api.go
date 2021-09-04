package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simple-go-auth/pkg/auth"
)

type Todo struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func CreateTodo(c *gin.Context) {
	var td Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	accessDetails := c.MustGet(auth.ContextAccessDetailsKey).(*auth.AccessDetails)
	userId, err := auth.GetAuthService().FetchAuthUserId(c, accessDetails.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	td.UserID = userId

	// return Todo struct
	c.JSON(http.StatusCreated, td)
}
func GetTodo(c *gin.Context) {
	accessDetails := c.MustGet(auth.ContextAccessDetailsKey).(*auth.AccessDetails)

	userId, err := auth.GetAuthService().FetchAuthUserId(c, accessDetails.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, Todo{
		UserID: userId,
		Title:  "Return from getting todo",
		Body:   "Return from getting todo for testing",
	})
}
