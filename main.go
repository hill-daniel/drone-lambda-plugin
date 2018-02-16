package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"os"
)

const defaultRegion = "us-east-1"

func main() {
	currentSession, err := createSession()
	if err != nil {
		log.Fatal("failed to create session", err)
	}
	svc := lambda.New(currentSession)

	result, err := svc.UpdateFunctionCode(createFunctionInput())

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatal(aerr.Code(), aerr.Error())
		} else {
			log.Fatal("failed to update function", err.Error())
		}
	}
	fmt.Println(result)
}

func createFunctionInput() *lambda.UpdateFunctionCodeInput {
	input := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(os.Getenv("PLUGIN_FUNCTION_NAME")),
		Publish:      aws.Bool(true),
		S3Bucket:     aws.String(os.Getenv("PLUGIN_S3_BUCKET")),
		S3Key:        aws.String(os.Getenv("PLUGIN_FILE_NAME")),
	}
	return input
}

func createSession() (*session.Session, error) {
	region := getRegion()
	currentSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	return currentSession, err
}

func getRegion() string {
	region, ok := os.LookupEnv("PLUGIN_FUNCTION_REGION")
	if !ok {
		region = defaultRegion
	}
	return region
}
