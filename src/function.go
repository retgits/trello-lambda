//go:generate go run ../../../TIBCOSoftware/flogo-lib/flogo/gen/gen.go $GOPATH

/*
Package main is the main executable of the serverless function. It will create a new
Trello card for each invocation of this service. To do so it requires access to Trello
using an appkey and an apptoken. Details on how to get thosr can be found in the
[Trello API documentation](https://trello.readme.io/docs/get-started).
*/
package main

// The imports
import (
	"errors"
	"fmt"
	"os"
	"reflect"

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

// Invoke is executed every time a new Lambda event is received.
// It takes the payload event (you can see an example in the event.json file) and
// returns a map[string]interface{} representing a JSON payload and an optional error
func Invoke(eventPayload interface{}) (map[string]interface{}, error) {
	trelloEvent := &TrelloEvent{}
	err := fillStruct(eventPayload.(map[string]interface{}), trelloEvent)
	if err != nil {
		return nil, err
	}

	in := map[string]interface{}{"token": accessToken, "appkey": appKey, "list": listID, "position": "bottom", "title": trelloEvent.Title, "description": trelloEvent.Description}
	out, err := flogo.EvalActivity(&trellocard.MyActivity{}, in)
	if err != nil {
		logger.Infof("Error while creating Trello card: %s", err.Error())
		return nil, err
	}

	logger.Infof("Trello card [%s] created successfully", trelloEvent.Title)
	return map[string]interface{}{"data": out["result"].Value(), "status": 200}, nil
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
