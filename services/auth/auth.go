package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"simple-go-auth/services/redis"
)

type AccessDetails struct {
	TokenUuid string
	Username  string
	UserId    string
}

type AuthInterface interface {
	CreateAuth(context.Context, string, *TokenDetails) error
	FetchAuthUserId(context.Context, string) (string, error)
	DeleteRefresh(context.Context, string) error
	DeleteTokens(context.Context, *AccessDetails) error
}

type AuthService struct {
}

var authServiceInstance *AuthService
var authServiceInstanceOnce sync.Once

func GetAuthService() *AuthService {
	authServiceInstanceOnce.Do(func() {
		authServiceInstance = &AuthService{}
	})
	return authServiceInstance
}

// Save token metadata to Redis
func (s *AuthService) CreateAuth(ctx context.Context, userId string, tokenDetails *TokenDetails) error {
	// converting Unix to UTC(to Time object)
	at := time.Unix(tokenDetails.AtExpires, 0)
	rt := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	atCreated, err := redis.GetRedisClient().Set(ctx, tokenDetails.TokenUuid, userId, at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := redis.GetRedisClient().Set(ctx, tokenDetails.RefreshUuid, userId, rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("failed to create the auth: no record inserted to redis")
	}
	return nil
}

// Check the metadata saved
func (s *AuthService) FetchAuthUserId(ctx context.Context, tokenUuid string) (string, error) {
	userid, err := redis.GetRedisClient().Get(ctx, tokenUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

// Once a user row in the token table
func (s *AuthService) DeleteTokens(ctx context.Context, authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%s", authD.TokenUuid, authD.UserId)
	//delete access token
	deletedAt, err := redis.GetRedisClient().Del(ctx, authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := redis.GetRedisClient().Del(ctx, refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("Can't delete metadata from Redis")
	}
	return nil
}

func (s *AuthService) DeleteRefresh(ctx context.Context, refreshUuid string) error {
	//delete refresh token
	deleted, err := redis.GetRedisClient().Del(ctx, refreshUuid).Result()
	if err != nil {
		return err
	} else if deleted == 0 {
		return errors.New("token expired")
	}
	return nil
}
