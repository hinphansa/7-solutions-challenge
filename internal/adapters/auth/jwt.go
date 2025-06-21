package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JWTMaker struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(secret string, ttl time.Duration) *JWTMaker {
	return &JWTMaker{secret: []byte(secret), ttl: ttl}
}

func (j *JWTMaker) NewToken(id bson.ObjectID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub": id.Hex(),
		"eml": email,
		"exp": time.Now().Add(j.ttl).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}
