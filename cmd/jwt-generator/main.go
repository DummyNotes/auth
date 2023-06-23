package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matelang/jwt-go-aws-kms/v2/jwtkms"
)

const keyID = ""

func main() {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	now := time.Now()
	jwtToken := jwt.NewWithClaims(jwtkms.SigningMethodECDSA512, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	})

	kmsConfig := jwtkms.NewKMSConfig(kms.NewFromConfig(awsConfig), keyID, false)

	str, err := jwtToken.SignedString(kmsConfig.WithContext(context.Background()))
	if err != nil {
		log.Fatalf("can not sign JWT %s", err)
	}

	fmt.Println(str)

	log.Printf("Signed JWT %s\n", str)
}
