package api

import (
	"net/http"

	"simple-go-auth/pkg/auth"
	"simple-go-auth/pkg/users"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
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

	ts, err := auth.GetTokenService().CreateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := auth.GetAuthService().CreateAuth(c, user.ID, ts); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func Logout(c *gin.Context) {
	// If metadata is passed and the tokens valid, delete them from the redis store
	accessDetails := c.MustGet(auth.ContextAccessDetailsKey).(*auth.AccessDetails)
	if accessDetails != nil {
		if err := auth.GetAuthService().DeleteTokens(c, accessDetails); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func Refresh(c *gin.Context) {
	tokenMap := map[string]string{}
	if err := c.ShouldBindJSON(&tokenMap); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	// verify the token
	token, err := auth.GetTokenService().VerifyTokenRefreshToken(tokenMap["refresh_token"])
	if err != nil {
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
		if err := auth.GetAuthService().DeleteRefresh(c, refreshUuid); err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		// Create new pairs of refresh and access tokens
		ts, err := auth.GetTokenService().CreateToken(userId, username)
		if err != nil {
			c.JSON(http.StatusForbidden, err.Error())
			return
		}
		// Save the tokens metadata to redis
		if err := auth.GetAuthService().CreateAuth(c, userId, ts); err != nil {
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
