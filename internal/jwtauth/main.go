package jwtauth

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matelang/jwt-go-aws-kms/v2/jwtkms"
)

type Claim struct {
	UserID string `json:"userid"`
	jwt.RegisteredClaims
}

func Generate(awsConfig aws.Config, keyID string) (jwtToken string, error error) {
	now := time.Now()

	claim := &Claim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID: "123",
	}

	token := jwt.NewWithClaims(jwtkms.SigningMethodECDSA512, claim)

	kmsConfig := jwtkms.NewKMSConfig(kms.NewFromConfig(awsConfig), keyID, false)

	str, err := token.SignedString(kmsConfig.WithContext(context.Background()))
	if err != nil {
		return "", err
	}

	return str, nil
}

func Validate(awsConfig aws.Config, keyID string, jwtToken string) (*jwt.RegisteredClaims, error) {
	kmsConfig := jwtkms.NewKMSConfig(kms.NewFromConfig(awsConfig), keyID, false)

	claims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(jwtToken,
		&claims, func(token *jwt.Token) (interface{}, error) {
			return kmsConfig, nil
		})
	if err != nil {
		return &claims, err
	}

	return &claims, nil
}
