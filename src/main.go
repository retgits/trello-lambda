//go:generate go run ../../../../TIBCOSoftware/flogo-lib/flogo/gen/gen.go $GOPATH

/*
Package main is the main executable of the serverless function. It will create a new
Trello card for each invocation of this service. To do so it requires access to Trello
using an appkey and an apptoken. Details on how to get thosr can be found in the
[Trello API documentation](https://trello.readme.io/docs/get-started).
*/
package main

// The imports
import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/aws/aws-lambda-go/lambda"
)

// LambdaEvent is the outer structure of the events that are received by this function
type LambdaEvent struct {
	EventVersion string
	EventSource  string
	Event        interface{}
}

// The handler function is executed every time that a new Lambda event is received.
// It takes a JSON payload (you can see an example in the event.json file) and only
// returns an error if the something went wrong.
func handler(request LambdaEvent) error {

	_, err := Invoke(request.Event)
	if err != nil {
		logger.Infof("Error while creating Trello card: %s", err.Error())
		return err
	}

	// Return no error
	return nil
}

// The main method is executed by AWS Lambda and points to the handler
func main() {
	lambda.Start(handler)
}
