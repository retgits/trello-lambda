// Package main is a the five minute journal function of Trello-Lambda.
// The app creates a new card containing the elements of the Five Minute
// Journal every day determined by the schedule set during deployment.
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

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

type Qod struct {
	Success  *Success  `json:"success,omitempty"`
	Contents *Contents `json:"contents,omitempty"`
}

type Contents struct {
	Quotes    []Quote `json:"quotes"`
	Copyright *string `json:"copyright,omitempty"`
}

type Quote struct {
	Quote      *string  `json:"quote,omitempty"`
	Length     *string  `json:"length,omitempty"`
	Author     *string  `json:"author,omitempty"`
	Tags       []string `json:"tags"`
	Category   *string  `json:"category,omitempty"`
	Date       *string  `json:"date,omitempty"`
	Permalink  *string  `json:"permalink,omitempty"`
	Title      *string  `json:"title,omitempty"`
	Background *string  `json:"background,omitempty"`
	ID         *string  `json:"id,omitempty"`
}

type Success struct {
	Total *int64 `json:"total,omitempty"`
}

// qodURL is the URL for the Quote of the Day API
const qodURL = "https://quotes.rest/qod"

const fiveMinJournalTemplate = `
{{ .Quote }}
{{ .Author }}

# Beginning of the day

## I am grateful for
_Name three things that you’re really grateful for._
1.
2.
3.

## What would make today great
_Name three things that if they would happen today, would make today really awesome_
1.
2.
3.

## Daily affirmations. I am...
_Describe two feelings you have about yourself._
<free text, two things like relaxed unangered/unflustered, calm and satisfied, happy about what I’ve done, unrushed, complete as is>

# End of the day

## Three amazing things that happened today
_Name three things that happened today that were really awesome_
1.
2.
3.

## How could I have made today better
_What could you have done that would have made today even better_
<free text>
`

var c config

func handler(request events.CloudWatchEvent) error {
	// Get configuration set using environment variables
	err := envconfig.Process("", &c)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", qodURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	q, err := unmarshalQod(body)
	if err != nil {
		return err
	}

	data := struct {
		Quote  string
		Author string
	}{
		Quote:  *q.Contents.Quotes[0].Quote,
		Author: *q.Contents.Quotes[0].Author,
	}

	cb, err := parseTemplate(fiveMinJournalTemplate, data)
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

	// Prepare the current date and the deadline date
	today := time.Now()
	deadline := time.Date(today.Year(), today.Month(), today.Day(), 17, 0, 0, 0, time.UTC)

	// Create a card
	card := &trello.Card{
		Name:   fmt.Sprintf("Five Minute Journal for %s", today.Format("2006-01-02")),
		Desc:   cb,
		Pos:    1,
		IDList: c.TrelloListID,
		Due:    &deadline,
	}

	return tc.CreateCard(card, trello.Defaults())
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

func unmarshalQod(data []byte) (Qod, error) {
	var r Qod
	err := json.Unmarshal(data, &r)
	return r, err
}

func parseTemplate(tpl string, data interface{}) (string, error) {
	var tplbuf bytes.Buffer

	parsedTpl, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return "", err
	}

	if err := parsedTpl.Execute(&tplbuf, data); err != nil {
		return "", err
	}

	return tplbuf.String(), nil
}

// The main method is executed by AWS Lambda and points to the handler
func main() {
	lambda.Start(handler)
}
