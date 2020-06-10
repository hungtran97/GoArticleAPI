package ulti

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	db "go-article/databases"
	"go-article/models"
)

func CheckLogin(username, password string) (models.User, error) {
	DB := db.ConnectDB()

	var user models.User
	DB.Where("username = ?", username).First(&user)

	err := models.CheckPasswordHash(user.PasswordHash, password)

	log.Println(user)

	return user, err
}

func CheckJwt(r *http.Request) (string, error) {
	var username string
	// Get jwt
	bearerToken := r.Header.Get("Authorization")
	tokenString := strings.Split(bearerToken, " ")[1]

	// Get Username from jwt
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_JWT")), nil
	})

	if err != nil {
		return username, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username = fmt.Sprintf("%v", claims["username"])
	}

	return username, nil
}

func CreateJwt(user models.User) (string, error) {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}
