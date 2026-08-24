package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleUser      Role = "USER"
	RoleAdmin     Role = "ADMIN"
	RoleModerator Role = "MODERATOR"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
	Role         Role
}

func GenerateToken(email, username string, userId primitive.ObjectID, addedTime time.Duration, role Role) (*TokenDetails, error) {

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(addedTime).Unix()
	td.AccessUUID = uuid.New().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 30).Unix()
	td.RefreshUUID = uuid.New().String()
	td.Role = role

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is not set")
	}

	// Create Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true // Fixing the typo here
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userId.Hex() // Convert ObjectID to string
	atClaims["exp"] = td.AtExpires
	atClaims["role"] = td.Role
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	td.AccessToken = accessToken
	fmt.Print("THIS IS ----->", accessToken)

	// Create Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userId.Hex() // Convert ObjectID to string
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(jwtSecret)) // Convert secret to []byte
	if err != nil {
		return nil, err
	}
	td.RefreshToken = refreshToken

	return td, nil
}

func VerifyToken(token string) (primitive.ObjectID, Role, error) {

	parseToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return primitive.ObjectID{}, "", errors.New("could not parse token")
	}

	if !parseToken.Valid {
		return primitive.ObjectID{}, "", errors.New("invalid token")

	}

	claims, ok := parseToken.Claims.(jwt.MapClaims)

	if !ok {
		return primitive.ObjectID{}, "", errors.New("invalid token claims")
	}

	userIdStr, ok := claims["user_id"].(string)

	if !ok {
		return primitive.ObjectID{}, "", errors.New("userId claim could not be parsed")
	}
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		return primitive.ObjectID{}, "", errors.New("userId claim is not a valid ObjectID")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return primitive.ObjectID{}, "", errors.New("role claim could not be parsed")
	}

	return userId, Role(role), nil
}

func VerifyRefreshToken(refreshTokenString string) (string, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return os.Getenv("JWT_SECRET_KEY"), nil
	})

	if err != nil {
		return "", err
	}

	// Extract token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string)
		if !ok {
			return "", fmt.Errorf("refresh UUID not found")
		}
		// Check if token is expired
		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			return "", fmt.Errorf("refresh token expired")
		}
		return refreshUUID, nil
	}

	return "", fmt.Errorf("invalid token")
}
