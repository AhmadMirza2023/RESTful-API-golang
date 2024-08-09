package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AhmadMirza2023/ecom/config"
	"github.com/AhmadMirza2023/ecom/types"
	"github.com/AhmadMirza2023/ecom/utils"
	"github.com/golang-jwt/jwt/v5"
)

var UserKey = &contextKey{"userID"}

type contextKey struct {
	name string
}

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the user request
		tokenString := getTokenFromRequest(r)
		log.Printf("Received token: %s", tokenString)

		// validate the JWT
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Print("invalid token")
			permissionDenied(w)
			return
		}

		// if it is we need to fetch the userID from the DB (id from the token)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Print("failed to parse token claims")
			permissionDenied(w)
			return
		}

		userIDStr, ok := claims["userID"].(string)
		if !ok {
			log.Print("userID claim is missing or not a string")
			permissionDenied(w)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(w)
			return
		}

		user, err := store.GetUserById(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// set context "userID" to the user ID
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" && strings.HasPrefix(tokenAuth, "Bearer ") {
		return strings.TrimPrefix(tokenAuth, "Bearer ")
	}
	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}
