package appjwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JwtUserInfo struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Exp      int64  `json:"exp"`
}

func CreateToken(user JwtUserInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":      time.Now().Unix(),
		"issuer":   "habitask",
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (JwtUserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return JwtUserInfo{}, err
	}
	if !token.Valid {
		return JwtUserInfo{}, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return JwtUserInfo{}, err
	}
	userInfo := JwtUserInfo{
		UserID:   claims["user_id"].(string),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
		Exp:      claims["expiresIn"].(int64),
	}
	return userInfo, nil
}
