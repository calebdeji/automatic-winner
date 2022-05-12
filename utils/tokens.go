package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var token_secret string = os.Getenv("TOKEN_KEY")

type SignedDetails struct {
	Email string
	Uid   string
	jwt.StandardClaims
}

func GenerateToken(claims *SignedDetails) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(token_secret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateNewTokenAndRefreshToken(email string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email: email,
		Uid:   uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := GenerateToken(claims)
	refreshToken, err := GenerateToken(refreshClaims)

	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, message string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(token_secret), nil
	})

	if err != nil {
		message = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		message = fmt.Sprintf("token is expired")
		message = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		message = fmt.Sprintf("token is expired")
		message = err.Error()
		return
	}

	return claims, message
}

func GenerateRandomOTP(length int) string {
	otp := ""

	for i := 0; i < length; i++ {
		value := rand.Intn(10)
		otp = otp + fmt.Sprint(value)
	}

	return otp
}
