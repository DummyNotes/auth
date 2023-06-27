package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dummynotes/notes/internal/jwtauth"
)

var KmsJwtKeyID = os.Getenv("KMS_JWT_KEY")

func main() {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	jwtToken, err := jwtauth.Generate(awsConfig, KmsJwtKeyID)
	if err != nil {
		log.Printf("can not sign JWT %s", err)
	}

	fmt.Println(jwtToken)

	log.Printf("Signed JWT %s\n", jwtToken)
}
