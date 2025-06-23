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

func (j *JWTMaker) Generate(id bson.ObjectID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub": id.Hex(),
		"eml": email,
		"exp": time.Now().Add(j.ttl).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j *JWTMaker) Verify(token string) (bson.ObjectID, bool, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil {
		return bson.ObjectID{}, false, err
	}
	id, err := bson.ObjectIDFromHex(claims["sub"].(string))
	if err != nil {
		return bson.ObjectID{}, false, err
	}
	return id, true, nil
}
