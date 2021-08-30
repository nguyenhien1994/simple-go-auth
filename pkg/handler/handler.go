package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"simple-go-auth/pkg/auth"
	"simple-go-auth/pkg/users"
)

// handler struct
type handler struct {
	authService auth.AuthInterface
	token       auth.TokenInterface
}

type Todo struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func NewHandlers(authService auth.AuthInterface, token auth.TokenInterface) *handler {
	return &handler{authService, token}
}

func (h *handler) Login(c *gin.Context) {
	var u users.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	user, err := users.FindUserByUsername(u.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	// compare the user from the request with sample user defined above
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	ts, err := h.token.CreateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := h.authService.CreateAuth(c, user.ID, ts); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *handler) Logout(c *gin.Context) {
	// If metadata is passed and the tokens valid, delete them from the redis store
	metadata, _ := h.token.ExtractTokenMetadata(c.Request)
	if metadata != nil {
		if err := h.authService.DeleteTokens(c, metadata); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func (h *handler) CreateTodo(c *gin.Context) {
	var td Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	metadata, err := h.token.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userId, err := h.authService.FetchAuthUserId(c, metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	td.UserID = userId

	// return Todo struct
	c.JSON(http.StatusCreated, td)
}
func (h *handler) GetTodo(c *gin.Context) {
	metadata, err := h.token.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := h.authService.FetchAuthUserId(c, metadata.TokenUuid)
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

func (h *handler) Refresh(c *gin.Context) {
	tokenMap := map[string]string{}
	if err := c.ShouldBindJSON(&tokenMap); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	// verify the token
	token, err := h.token.VerifyTokenRefreshToken(tokenMap["refresh_token"])
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid refresh token")
		return
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "unauthorized: missing refresh_uuid")
			return
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "unauthorized: missing user_id")
			return
		}
		username, ok := claims["user_name"].(string)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "unauthorized: missing user_name")
			return
		}
		// Check refresh token in Redis and delete the previous refresh token
		if err := h.authService.DeleteRefresh(c, refreshUuid); err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		// Create new pairs of refresh and access tokens
		ts, err := h.token.CreateToken(userId, username)
		if err != nil {
			c.JSON(http.StatusForbidden, err.Error())
			return
		}
		// Save the tokens metadata to redis
		if err := h.authService.CreateAuth(c, userId, ts); err != nil {
			c.JSON(http.StatusForbidden, err.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh token expired")
	}
}
