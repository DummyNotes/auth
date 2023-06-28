package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dummynotes/notes/internal/jwtauth"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var KmsJwtKeyID = os.Getenv("KMS_JWT_KEY_ID")

func main() {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	jwtToken, err := jwtauth.Generate(awsConfig, KmsJwtKeyID)
	if err != nil {
		log.Errorf("can not sign JWT %s", err)
	}

	fmt.Println(jwtToken)

	log.Infof("Signed JWT %s\n", jwtToken)
}
