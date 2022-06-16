package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"medods/internal/adapters/database"
	"medods/internal/model"
	"medods/pkg"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type service struct {
	storage database.Storage
}

type AuthServiceInterface interface {
	CreateToken(userGUID string) (*model.Jwt, error)
	UpdateToken(tokens *model.Jwt) (*model.Jwt, error)
}

type AuthClaims struct {
	*jwt.StandardClaims
	User *model.UserToken
}

func NewService(storage database.Storage) AuthServiceInterface {
	return &service{
		storage: storage,
	}
}

func (s *service) CreateToken(userGUID string) (*model.Jwt, error) {
	tokens, userRefreshToken, err := s.createPairOfTokens(userGUID)
	if err != nil {
		log.Printf("don't create pair of tokens: %v\n", err)
		return nil, err
	}
	s.storage.CreateOne(userRefreshToken)
	return tokens, nil
}

func (s *service) UpdateToken(tokens *model.Jwt) (*model.Jwt, error) {
	_, valid := s.checkAccessToken(tokens.AccsessToken)
	if !valid {
		log.Printf("this access token is not valid:--> %v\n", valid)
		return nil, errors.New("this access token is not valid")
	}
	bindTokens := tokens.AccsessToken[len(tokens.AccsessToken)-6:]
	oldRefresh, err := s.storage.GetOne(tokens.UserGUID, bindTokens)
	if err != nil {
		log.Printf("Did't find the second pair of token: %v\n", err)
		return nil, err
	}
	isValidRefresh := pkg.CompareHashAndData(tokens.RefreshToken, oldRefresh.RefreshToken)
	if !isValidRefresh {
		log.Println("this refresh token is not valid")
		return nil, errors.New("this refresh token is not valid")
	}
	newJwtTokens, newRefreshToken, err := s.createPairOfTokens(tokens.UserGUID)
	if err != nil {
		log.Printf("don't create pair of tokens: %v\n", err)
		return nil, err
	}
	s.storage.UpdateOne(oldRefresh, newRefreshToken)
	return newJwtTokens, nil
}

func (s *service) createPairOfTokens(userGUID string) (*model.Jwt, *model.UserToken, error) {
	accessToken, err := s.createAccessToken(userGUID)
	if err != nil {
		log.Printf("don't create access token: %v\n", err)
		return nil, nil, err
	}
	refreshToken := s.createRefreshToken()
	refreshTokenHash, err := pkg.GenerateHash(refreshToken)
	if err != nil {
		log.Printf("don't create refresh token hash: %v\n", err)
		return nil, nil, err
	}

	// bind an access token with a refresh token
	bindTokens := accessToken[len(accessToken)-6:]

	tokensForClient := &model.Jwt{
		UserGUID:     userGUID,
		AccsessToken: accessToken,
		RefreshToken: refreshToken,
	}
	tokensForDB := &model.UserToken{
		UserGUID:     userGUID,
		RefreshToken: refreshTokenHash,
		BindTokens:   bindTokens,
	}
	return tokensForClient, tokensForDB, nil
}

func (s *service) createAccessToken(userGUID string) (string, error) {
	accessTokenClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		Id:        userGUID,
		IssuedAt:  time.Now().Add(time.Second * 5).Unix(),
		Subject:   "user",
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	token, err := at.SignedString([]byte(os.Getenv("access_key")))
	if err != nil {
		return "", err
	}
	fmt.Println("acsess token value:-->", token)

	return token, nil
}

func (s *service) createRefreshToken() string {
	refreshToken := uuid.New()
	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken.String()))
	return refreshTokenBase64
}

func (s *service) checkAccessToken(accessToken string) (*model.UserToken, bool) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// return []byte("medodskey"), nil
		return []byte(os.Getenv("access_key")), nil
	})
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			log.Println("time EXPIRED")
			return nil, true
		}
		return nil, false
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, true
	}
	return nil, false
}
