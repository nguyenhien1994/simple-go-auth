package auth

import (
	"context"
	"errors"
	"time"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type AccessDetails struct {
	TokenUuid string
	UserId    string
}

type AuthInterface interface {
	CreateAuth(string, *TokenDetails) error
	FetchAuth(string) (string, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

type service struct {
	client *redis.Client
}

func NewAuthService(client *redis.Client) *service {
	return &service{client: client}
}

//Save token metadata to Redis
func (s *service) CreateAuth(userId string, tokenDetails *TokenDetails) error {
	//converting Unix to UTC(to Time object)
	at := time.Unix(tokenDetails.AtExpires, 0)
	rt := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	atCreated, err := s.client.Set(context.Background(), tokenDetails.TokenUuid, userId, at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := s.client.Set(context.Background(), tokenDetails.RefreshUuid, userId, rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("failed to create the auth: no record inserted to redis")
	}
	return nil
}

//Check the metadata saved
func (s *service) FetchAuth(tokenUuid string) (string, error) {
	userid, err := s.client.Get(context.Background(), tokenUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

//Once a user row in the token table
func (s *service) DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%s", authD.TokenUuid, authD.UserId)
	//delete access token
	deletedAt, err := s.client.Del(context.Background(), authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := s.client.Del(context.Background(), refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (s *service) DeleteRefresh(refreshUuid string) error {
	//delete refresh token
	deleted, err := s.client.Del(context.Background(), refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
