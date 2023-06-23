package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matelang/jwt-go-aws-kms/v2/jwtkms"
	"log"
)

const keyID = ""

func HandleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (*events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	kmsConfig := jwtkms.NewKMSConfig(kms.NewFromConfig(awsConfig), keyID, false)

	claims := jwt.RegisteredClaims{}

	_, err = jwt.ParseWithClaims(request.Headers["authorization"],
		&claims, func(token *jwt.Token) (interface{}, error) {
			return kmsConfig, nil
		})
	if err != nil {
		log.Printf("can not parse/verify token %s", err)

		return simpleResponse(false), nil
	}

	log.Printf("Parsed and validated token with claims %v", claims)

	return simpleResponse(true), nil

}

func simpleResponse(isAuthorized bool) *events.APIGatewayV2CustomAuthorizerSimpleResponse {
	return &events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: isAuthorized,
		Context:      make(map[string]interface{}),
	}
}

func main() {
	lambda.Start(HandleRequest)
}
