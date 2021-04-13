package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Example: https://github.com/dgrijalva/jwt-go/blob/master/hmac_example_test.go

func CreateToken(user_id uint32) (string, error) {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 744).Unix() // 31 Days

	// Set what cryptographic algorithms used + claims -> header and payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signed -> Signature
	return token.SignedString([]byte(os.Getenv("APP_API_SECRET")))
}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)

	//This part follow the example link above
	// https://tour.golang.org/methods/15 -> Method.(*jwt.SigningMethodHMAC)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("APP_API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

func ExtractTokenID(r *http.Request) (uint32, error) {

	tokenString := ExtractToken(r)
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	// Check if Claims is MapClaims and return value inside Claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}

func ExtractToken(r *http.Request) string {
	// Token in URL
	keys := r.URL.Query()
	token := keys.Get("token")

	if token != "" {
		return token
	}

	// Token in Header
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

//Pretty display the claims licely in the terminal
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
