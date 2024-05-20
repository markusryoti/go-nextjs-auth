package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/markusryoti/next-js-go-auth/store"
	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 10

var mySigningKey = []byte("verysecret")

type SessionResponse struct {
	SessionId string `json:"sessionId"`
	ExpiresAt int64  `json:"expiresAt"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type MyCustomClaims struct {
	Session string `json:"session"`
	jwt.RegisteredClaims
}

func RegisterUser(email, password string) (string, error) {

	existing, err := store.UserStore.GetUserByEmail(email)
	if err != nil && !errors.Is(err, store.ErrUserNotFound) {
		return "", err
	}

	if existing != nil {
		return "", errors.New("user already found")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return "", err
	}

	hashedPassword := string(hashedBytes)

	user, err := store.UserStore.SaveUser(email, hashedPassword)
	if err != nil {
		return "", err
	}

	sessionId := uuid.NewString()
	expiresAt := time.Now().Add(time.Hour * 24 * 30)

	// TODO
	// Add expiry time
	store.SessionStore.AddSession(sessionId, user)

	token, err := generateAccessToken(sessionId, expiresAt)
	if err != nil {
		return "", err
	}

	return token, err
}

func LoginUser(email, password string) (TokenResponse, error) {
	var res TokenResponse

	user, err := store.UserStore.GetUserByEmail(email)
	if err != nil {
		return res, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return res, errors.New("passwords don't match")
	}

	sessionId := uuid.NewString()
	expiresAt := time.Now().Add(time.Hour * 24 * 30)

	store.SessionStore.AddSession(sessionId, user)

	accessToken, err := generateAccessToken(sessionId, expiresAt)
	if err != nil {
		return res, err
	}

	res.AccessToken = accessToken

	return res, nil
}

func ValidateToken(token string) (MyCustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return MyCustomClaims{}, err
	}

	claims, ok := parsedToken.Claims.(*MyCustomClaims)
	if !ok {
		return MyCustomClaims{}, errors.New("couldn't parse custom claims")
	}

	// expired := claims.ExpiresAt.Before(time.Now())
	// if expired {
	// 	return MyCustomClaims{}, errors.New("token expired")
	// }

	return *claims, nil
}

func generateAccessToken(sessionId string, expiresAt time.Time) (string, error) {
	claims := createClaims(sessionId, expiresAt)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	return ss, err
}

func createClaims(sessionId string, expiresAt time.Time) MyCustomClaims {
	claims := MyCustomClaims{
		sessionId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	return claims
}
