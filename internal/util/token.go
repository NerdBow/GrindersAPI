package util

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	CanNotAccessJWTSecretErr = errors.New("Unable to get JWT secret from env")
)

func CreateToken(claimMap jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimMap)

	jwtSecret := os.Getenv("JWTSECRET")

	if jwtSecret == "" {
		log.Println("CRITICAL: JWT secret could not be found in the .env file!")
		return "", CanNotAccessJWTSecretErr
	}

	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", errors.New("Could not sign")
	}

	return signedToken, nil

}

func GetClaimsFromToken(token string) (jwt.MapClaims, error) {
	key := func(token *jwt.Token) (any, error) {
		secret := os.Getenv("JWTSECRET")
		if secret == "" {
			return nil, CanNotAccessJWTSecretErr
		}
		return []byte(secret), nil
	}

	parsedToken, err := jwt.Parse(token, key, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		log.Println("Token has no claims.")
		return nil, err
	}
	return claims, nil
}

func CheckTokenExpiration(claims jwt.MapClaims) (bool, error) {
	exp, err := claims.GetExpirationTime()

	if err != nil {
		return false, err
	}

	if time.Now().After(exp.Time) {
		return false, jwt.ErrTokenExpired
	}

	return true, nil
}
