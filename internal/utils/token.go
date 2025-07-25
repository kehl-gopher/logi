package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, secret string, duration time.Duration) (string, time.Time, error) {
	expiration := time.Now().Add(duration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiration.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiration, nil
}

func ValidateToken(tokenStr string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token provided")
}

func GetUserIDFromToken(tokenStr string, secretKey string) (string, error) {
	claims, err := ValidateToken(tokenStr, secretKey)
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	return userID, nil
}

func ParseTime(input string) (time.Duration, error) {
	inputL := strings.ToLower(input)
	if strings.HasSuffix(inputL, "d") {
		d := strings.TrimSuffix(inputL, "d")
		toInt, err := strconv.Atoi(d)
		if err != nil {
			return 0, err
		}
		return time.Duration(toInt) * 24 * time.Hour, nil
	}
	if strings.HasSuffix(inputL, "hr") ||
		strings.HasSuffix(inputL, "h") ||
		strings.HasSuffix(inputL, "m") ||
		strings.HasSuffix(inputL, "s") {
		t, err := time.ParseDuration(input)
		if err != nil {
			return 0, nil
		}
		return t, nil
	}
	return 0, fmt.Errorf("invalid time unit provided should be support d -> day, hr -> hour, m -> minutes, s -> second")
}
