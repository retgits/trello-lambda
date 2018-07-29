//go:generate go run ../../../TIBCOSoftware/flogo-lib/flogo/gen/gen.go -shim $GOPATH

/*
Package main is the main executable of the serverless function. It will create a new
Trello card for each invocation of this service. To do so it requires access to Trello
using an appkey and an apptoken. Details on how to get thosr can be found in the
[Trello API documentation](https://trello.readme.io/docs/get-started).
*/
package main

// The imports
import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/lambda"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/flogo"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/retgits/flogo-components/activity/trellocard"
)

var (
	// Your Trello App token
	appKey = os.Getenv("appkey")

	// Your Trello App key
	accessToken = os.Getenv("apptoken")

	// The ID of the list you want to send the card to
	listID = os.Getenv("list")
)

// TrelloEvent is the structure for the data representing a TrelloCard
type TrelloEvent struct {
	Title       string
	Description string
}

// shimApp creates a Flogo app with the Lambda trigger and registers function handlers to be executed. The return value is
// a pointer to a Flogo app that is used by the shim_support file to create the engine.
func shimApp() *flogo.App {

	app := flogo.NewApp()

	trg := app.NewTrigger(&lambda.LambdaTrigger{}, nil)
	trg.NewFuncHandler(nil, RunActivities)

	return app
}

// RunActivities is executed every time that a new Lambda event is received.
// It takes a JSON payload (you can see an example in the event.json file) and only
// returns an error if the something went wrong.
func RunActivities(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	// Create a TrelloEvent from the incoming event
	trelloEvent := &TrelloEvent{}
	eventMap := inputs["evt"].Value().(map[string]interface{})
	err := fillStruct(eventMap["Event"].(map[string]interface{}), trelloEvent)
	if err != nil {
		return nil, err
	}

	// Run the Trello activity
	in := map[string]interface{}{"token": accessToken, "appkey": appKey, "list": listID, "position": "bottom", "title": trelloEvent.Title, "description": trelloEvent.Description}
	_, err = flogo.EvalActivity(&trellocard.MyActivity{}, in)
	if err != nil {
		logger.Infof("Error while creating Trello card: %s", err.Error())
		return nil, err
	}
	logger.Infof("Trello card [%s] created successfully", trelloEvent.Title)

	// Return nothing
	return nil, nil
}

// fillStruct maps the fields from a map[string]interface{} to a struct
func fillStruct(m map[string]interface{}, s interface{}) error {
	structValue := reflect.ValueOf(s).Elem()

	for name, value := range m {
		structFieldValue := structValue.FieldByName(name)

		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in obj", name)
		}

		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", name)
		}

		val := reflect.ValueOf(value)
		if structFieldValue.Type() != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)
	}
	return nil
}
