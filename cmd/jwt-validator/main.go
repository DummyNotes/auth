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

var KmsJwtKeyID = os.Getenv("KMS_JWT_KEY")

func HandleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (*events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	responseContext := make(map[string]interface{})

	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	claims, err := jwtauth.Validate(awsConfig, KmsJwtKeyID, request.Headers["authorization"])
	if err != nil {
		log.Errorf("can not parse/verify token %s", err)
		return nil, err
	}

	log.Infof("Parsed and validated token with claims %v", claims)

	responseContext["userid"] = "1233456456"

	return simpleResponse(true, responseContext), nil

}

func simpleResponse(isAuthorized bool, responseContext map[string]interface{}) *events.APIGatewayV2CustomAuthorizerSimpleResponse {
	return &events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: isAuthorized,
		Context:      responseContext,
	}
}

func main() {
	lambda.Start(HandleRequest)
}
