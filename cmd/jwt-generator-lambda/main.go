package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dummynotes/notes/internal/jwtauth"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var KmsJwtKeyID = os.Getenv("KMS_JWT_KEY_ID")

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       map[string]string `json:"body"`
}

func HandleRequest(ctx context.Context) (Response, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	jwtToken, err := jwtauth.Generate(awsConfig, KmsJwtKeyID)
	if err != nil {
		log.Errorf("can not sign JWT %s", err)
	}

	log.Infof("Signed JWT %s\n", jwtToken)

	return Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       map[string]string{"jwtToken": jwtToken},
		},
		nil
}

func main() {
	lambda.Start(HandleRequest)
}
