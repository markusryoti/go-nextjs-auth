package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base32"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markusryoti/next-js-go-auth/internal"
	"github.com/markusryoti/next-js-go-auth/internal/store"
	"golang.org/x/crypto/argon2"
)

const (
	SessionExpiryTime = 30 * 24 * time.Hour
	ArgonTime         = 2
	ArgonMemory       = 19 * 1024
	ArgonKeyLen       = 32
)

var mySigningKey = []byte("verysecret")

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type MyCustomClaims struct {
	Session string `json:"session"`
	jwt.RegisteredClaims
}

func RegisterUser(email, password string) (TokenResponse, error) {
	var res TokenResponse

	existing, err := store.UserStore.GetUserByEmail(email)
	if err != nil && !errors.Is(err, store.ErrUserNotFound) {
		return res, err
	}

	if existing != nil {
		return res, errors.New("user already found")
	}

	salt := generateSalt()

	hashedBytes := argon2.IDKey([]byte(password), salt, ArgonTime, ArgonMemory, 1, ArgonKeyLen)

	hashedPassword := string(hashedBytes)

	user, err := store.UserStore.SaveUser(email, hashedPassword, string(salt))
	if err != nil {
		return res, err
	}

	sessionId := genSessionId()
	expiryTime := time.Now().Add(SessionExpiryTime)

	store.SessionStore.AddSession(sessionId, user, expiryTime)

	token, err := generateAccessToken(sessionId, expiryTime)
	if err != nil {
		return res, err
	}

	res.AccessToken = token

	return res, err
}

func LoginUser(email, password string) (TokenResponse, error) {
	var res TokenResponse

	user, err := store.UserStore.GetUserByEmail(email)
	if err != nil {
		return res, err
	}

	hash := argon2.IDKey([]byte(password), []byte(user.Salt), ArgonTime, ArgonMemory, 1, ArgonKeyLen)

	if subtle.ConstantTimeCompare(hash, []byte(user.HashedPassword)) == 0 {
		return res, errors.New("passwords don't match")
	}

	sessionId := genSessionId()
	expiryTime := time.Now().Add(SessionExpiryTime)

	store.SessionStore.AddSession(sessionId, user, expiryTime)

	token, err := generateAccessToken(sessionId, expiryTime)
	if err != nil {
		return res, err
	}

	res.AccessToken = token

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

	return *claims, nil
}

type CurrentUserResponse struct {
	User        *internal.User `json:"user"`
	AccessToken *string        `json:"accessToken"`
}

func GetCurrentUser(sessionId string) (*CurrentUserResponse, error) {
	session, err := store.SessionStore.GetSession(sessionId)
	if err != nil {
		return nil, err
	}

	if session.Expired() {
		return nil, errors.New("session expired")
	}

	then := session.ExpiresAt.Add(-SessionExpiryTime / 2)

	var accessToken string

	if time.Now().After(then) {
		newExpiry := time.Now().Add(SessionExpiryTime)

		session.ExpiresAt = newExpiry
		store.SessionStore.UpdateExpiration(sessionId, newExpiry)

		accessToken, err = generateAccessToken(sessionId, newExpiry)
		if err != nil {
			return nil, err
		}
	}

	res := &CurrentUserResponse{
		User:        session.User,
		AccessToken: &accessToken,
	}

	return res, nil
}

func Logout(sessionId string) {
	store.SessionStore.RemoveSession(sessionId)
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

func genSessionId() string {
	bytes := make([]byte, 15)
	rand.Read(bytes)
	sessionId := base32.StdEncoding.EncodeToString(bytes)
	return sessionId
}

func generateSalt() []byte {
	bytes := make([]byte, 15)
	rand.Read(bytes)
	return bytes
}
