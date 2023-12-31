package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dummynotes/notes/internal/jwtauth"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var KmsJwtKeyID = os.Getenv("KMS_JWT_KEY_ID")

func HandleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	claims, err := jwtauth.Validate(awsConfig, KmsJwtKeyID, request.Headers["authorization"])
	if err != nil {
		log.Errorf("can not parse/verify token %s", err)
		log.Infof("authorization header: %s", request.Headers["authorization"])
		log.Infof("KMS_JWT_KEY_ID: %s", KmsJwtKeyID)

		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	log.Infof("Parsed and validated token with claims %v", claims)

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
