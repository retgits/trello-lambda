// Package main is a the weekly archiver function of Trello-Lambda.
// The app archives all cards in the "done" folder on a weekly basis
// determined by the schedule set during deployment.
package main

import (
	"encoding/base64"
	"fmt"

	"github.com/adlio/trello"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	AWSRegion      string `required:"true" split_words:"true" envconfig:"AWS_REGION"`
	TrelloAPIKey   string `required:"true" split_words:"true" envconfig:"API_KEY"`
	TrelloAppToken string `required:"true" split_words:"true" envconfig:"APP_TOKEN"`
	TrelloListID   string `required:"true" split_words:"true" envconfig:"LIST_ID"`
}

var c config

func handler(request events.CloudWatchEvent) error {
	// Get configuration set using environment variables
	err := envconfig.Process("", &c)
	if err != nil {
		return err
	}

	// Create a new AWS session
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(c.AWSRegion),
	}))

	// Create a new KMS session
	kmsSvc := kms.New(awsSession)

	// Update the API Key with the decoded value
	val, err := decodeString(kmsSvc, c.TrelloAPIKey)
	if err != nil {
		return err
	}
	c.TrelloAPIKey = val

	// Update the App Token with the decoded value
	val, err = decodeString(kmsSvc, c.TrelloAppToken)
	if err != nil {
		return err
	}
	c.TrelloAppToken = val

	// Create a connection to Trello
	tc := trello.NewClient(c.TrelloAPIKey, c.TrelloAppToken)

	// Get the list based on the list ID
	l, err := tc.GetList(c.TrelloListID, trello.Defaults())
	if err != nil {
		return err
	}

	// Get the cards on the list
	cs, err := l.GetCards(trello.Defaults())
	if err != nil {
		return err
	}

	for _, card := range cs {
		args := make(trello.Arguments)
		args["closed"] = "true"
		err := card.Update(args)
		if err != nil {
			fmt.Printf("error updating %s: %s", card.ID, err.Error())
		}
	}

	return nil
}

// decodeString uses AWS Key Management Service (AWS KMS) to decrypt environment variables.
// In order for this method to work, the function needs access to the kms:Decrypt capability.
func decodeString(kmsSvc *kms.KMS, payload string) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}
	out, err := kmsSvc.Decrypt(&kms.DecryptInput{
		CiphertextBlob: sDec,
	})
	if err != nil {
		return "", err
	}
	return string(out.Plaintext), nil
}

// The main method is executed by AWS Lambda and points to the handler
func main() {
	lambda.Start(handler)
}
