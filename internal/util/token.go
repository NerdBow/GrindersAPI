package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strconv"
	"time"
)

func createToken(user model.User) (string, error) {
	jwtExpDuration, err := strconv.Atoi(os.Getenv("JWTEXP"))

	if err != nil {
		log.Println("CRITICAL: JWTEXP could not be parsed from the .env file!")
		return "", errors.New("Can not get JWT EXP time from env")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.UserId,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * time.Duration(jwtExpDuration)).Unix(),
	})

	jwtSecret := os.Getenv("JWTSECRET")

	if jwtSecret == "" {
		log.Println("CRITICAL: JWT secret could not be found in the .env file!")
		return "", errors.New("Can not get secret")
	}

	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", errors.New("Could not sign")
	}

	return signedToken, nil

}

func getClaimsFromToken(token string) (jwt.MapClaims, error) {
	key := func(token *jwt.Token) (any, error) {
		secret := os.Getenv("JWTSECRET")
		if secret == "" {
			return nil, errors.New("Can not get secret")
		}
		return []byte(secret), nil
	}

	parsedToken, err := jwt.Parse(token, key, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		log.Println("Token has no claims.")
		return nil, err
	}
	return claims, nil
}

func checkTokenExpiration(claims jwt.MapClaims) (bool, error) {
	exp, err := claims.GetExpirationTime()

	if err != nil {
		return false, err
	}

	if time.Now().After(exp.Time) {
		return false, jwt.ErrTokenExpired
	}

	return true, nil
}

func getUserFromClaims(claims jwt.MapClaims) (model.User, error) {
	var user model.User

	ok, err := checkTokenExpiration(claims)

	if err != nil {
		return user, err
	}

	user.UserId, ok = claims["userId"].(int)

	if !ok {
		return user, jwt.ErrTokenMalformed
	}

	user.Username, ok = claims["username"].(string)

	if !ok {
		return user, jwt.ErrTokenMalformed
	}

	return user, nil
}
