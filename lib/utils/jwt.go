package utils

import (
	"errors"
	"flaver/globals"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type FlaverCliaims struct {
	jwt.StandardClaims

	Role string `json:"role"`
}

type JwtResult struct {
	Token     string
	ExpiredAt time.Time
}

func JwtTokenGenerator(userUid string, role string, expiredInHour int) (*JwtResult, error) {
	now := time.Now()
	claims := FlaverCliaims{
		jwt.StandardClaims{
			Audience:  userUid,
			ExpiresAt: now.Add(time.Duration(expiredInHour) * time.Hour).Unix(),
		},
		role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(globals.GetConfig().Jwt.Secret)
	if err != nil {
		return nil, errors.New("signed token string error")
	}
	return &JwtResult{
		Token: tokenString,
	}, nil
}

func ParseJWTToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &FlaverCliaims{}, func(token *jwt.Token) (interface{}, error) {
		return globals.GetConfig().Jwt.Secret, nil
	})
}

func ValidateJWTToken(token string) (*jwt.Token, error) {
	jwtToken, err := ParseJWTToken(token)
	// check if token is invalid
	if err != nil || !jwtToken.Valid {
		return nil, err
	}

	// check if token has format error
	_, ok := jwtToken.Claims.(*FlaverCliaims)
	if !ok {
		return nil, err
	}

	return jwtToken, nil
}

func GetUIDFromToken(jwtToken *jwt.Token) string {
	claims, _ := jwtToken.Claims.(*FlaverCliaims)

	return claims.Audience
}
