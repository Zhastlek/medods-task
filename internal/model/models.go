package model

import "github.com/golang-jwt/jwt"

type Jwt struct {
	*jwt.StandardClaims
	UserGUID     string `json:"user_guid"`
	AccsessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserToken struct {
	UserGUID     string `bson:"user_guid"`
	RefreshToken string `bson:"refresh_token"`
	BindTokens   string `bson:"bind_tokens"`
}
