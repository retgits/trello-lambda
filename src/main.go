/*
Package main is the main executable of the serverless function. It will create a new
Trello card for each invocation of this service. To do so it requires access to Trello
using an appkey and an apptoken. Details on how to get thosr can be found in the
[Trello API documentation](https://trello.readme.io/docs/get-started).
*/
package main

// The imports
import (
	"log"
	"os"

	"github.com/adlio/trello"
	"github.com/aws/aws-lambda-go/lambda"
)

// Variables that are set as Environment Variables
var (
	trelloAppKey      = os.Getenv("appkey")
	trelloAccessToken = os.Getenv("apptoken")
	trelloListID      = os.Getenv("defaultlist")
)

// The handler function is executed every time that a new Lambda event is received.
// It takes a JSON payload (you can see an example in the event.json file) and only
// returns an error if the something went wrong.
func handler(request map[string]interface{}) error {
	// Create a new Trello client
	trelloClient := trello.NewClient(trelloAppKey, trelloAccessToken)

	// Get the title and the description for the card
	trelloEvent := request["Trello"].(map[string]interface{})
	cardTitle := trelloEvent["Title"].(string)
	cardDescription := trelloEvent["Description"].(string)
	log.Printf("Got a new event: %s", trelloEvent)

	// Create an instance of a card
	card := trello.Card{
		Name:   cardTitle,
		Desc:   cardDescription,
		IDList: trelloListID,
	}

	// Create the card on the Trello board
	err := trelloClient.CreateCard(&card, trello.Defaults())

	if err != nil {
		log.Print(err)
		return err
	}

	// Move the card to the bottom of the list
	card.MoveToBottomOfList()

	// Return no error
	return nil
}

// The main method is executed by AWS Lambda and points to the handler
func main() {
	lambda.Start(handler)
}
