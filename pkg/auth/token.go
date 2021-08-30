package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type TokenService struct {
	accessSecret  string
	refreshSecret string
}

var tokenUtilsInstance *TokenService
var tokenUtilsInstanceOnce sync.Once

func GetTokenService() *TokenService {
	tokenUtilsInstanceOnce.Do(func() {
		tokenUtilsInstance = &TokenService{os.Getenv("ACCESS_SECRET"), os.Getenv("REFRESH_SECRET")}
	})
	return tokenUtilsInstance
}

type TokenInterface interface {
	CreateToken(userId string, username string) (*TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*AccessDetails, error)
	TokenValid(r *http.Request) error
	VerifyTokenRefreshToken(tokenStr string) (*jwt.Token, error)
}

func (t *TokenService) CreateToken(userId string, username string) (*TokenDetails, error) {
	td := generateNewTokenDetails(userId)

	var err error
	// Creating access token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.TokenUuid
	atClaims["user_id"] = userId
	atClaims["user_name"] = username
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(t.accessSecret))
	if err != nil {
		return nil, err
	}

	// Creating refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["user_name"] = username
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(t.refreshSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (t *TokenService) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := verifyAccessToken(r)
	if err != nil {
		return nil, err
	}
	acc, err := extract(token)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (t *TokenService) TokenValid(r *http.Request) error {
	token, err := verifyAccessToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func (t *TokenService) VerifyTokenRefreshToken(tokenStr string) (*jwt.Token, error) {
	// verify the token
	return verifyToken(tokenStr, t.refreshSecret)
}

// get the token from the request body
func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func extract(token *jwt.Token) (*AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(string)
		username, usernameOk := claims["user_name"].(string)
		if ok == false || userOk == false || usernameOk == false {
			return nil, errors.New("unauthorized")
		} else {
			return &AccessDetails{
				TokenUuid: accessUuid,
				UserId:    userId,
				Username:  username,
			}, nil
		}
	}
	return nil, errors.New("failed to extract token from request")
}

func generateNewTokenDetails(userId string) *TokenDetails {
	td := &TokenDetails{}
	// access token expires after 10 min
	td.AtExpires = time.Now().Add(time.Minute * 10).Unix()
	td.TokenUuid = uuid.NewV4().String()

	// refresh token expires after 24h
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + userId

	return td
}

func verifyToken(tokenStr string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func verifyAccessToken(r *http.Request) (*jwt.Token, error) {
	return verifyToken(extractToken(r), os.Getenv("ACCESS_SECRET"))
}
