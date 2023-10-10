package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTMaker struct{
	secretKey string
}
var minSecretKeySize = 32
func NewJWTMaker(secretKey string)(Maker, error){
	if len(secretKey) < minSecretKeySize{
		 return nil, fmt.Errorf("invalid key size: must be atleast %d characters", minSecretKeySize)
	}
	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func(j *JWTMaker)CreateToken(username string, duration time.Duration)(string, error){
	payload, err := NewPayload(username, duration)
	if err != nil{
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,payload)
	return token.SignedString([]byte(j.secretKey))
}

func(j *JWTMaker)VerifyToken(token string)(*Payload, error){
	keyFunc := func(token *jwt.Token) (interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok{
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}
	jwToken, err := jwt.ParseWithClaims(token,&Payload{}, keyFunc)
	if err != nil{
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner,ErrTokenExpired){
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	 payload, ok := jwToken.Claims.(*Payload)
	 if !ok{
		return nil, ErrInvalidToken
	 }
	 return payload, nil

}
