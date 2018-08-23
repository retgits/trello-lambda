//go:generate go run $GOPATH/src/github.com/TIBCOSoftware/flogo-lib/flogo/gen/gen.go $GOPATH

// Package main is the main executable of the serverless function. It will create a new
// Trello card for each invocation of this service. To do so it requires access to Trello
// using an appkey and an apptoken. Details on how to get thosr can be found in the
// [Trello API documentation](https://trello.readme.io/docs/get-started).
package main

// The ever important imports
import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/lambda"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/flogo"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/retgits/flogo-components/activity/awsssm"
	"github.com/retgits/flogo-components/activity/trellocard"
)

// Constants
const (
	// The name of the Trello App Key parameter in Amazon SSM
	appkey = "/trello/appkey"
	// The name of the Trello App Token parameter in Amazon SSM
	apptoken = "/trello/apptoken"
	// The name of the default Trello list parameter in Amazon SSM
	applist = "/trello/list"
)

// TrelloEvent is the structure for the data representing a TrelloCard
type TrelloEvent struct {
	Title       string
	Description string
}

// Init makes sure that everything is ready to go!
func init() {
	config.SetDefaultLogLevel("INFO")
	logger.SetLogLevel(logger.InfoLevel)

	app := shimApp()

	e, err := flogo.NewEngine(app)

	if err != nil {
		logger.Error(err)
		return
	}

	e.Init(true)
}

// shimApp is used to build a new Flogo app and register the Lambda trigger with the engine.
// The shimapp is used by the shim, which triggers the engine every time an event comes into Lambda.
func shimApp() *flogo.App {
	// Create a new Flogo app
	app := flogo.NewApp()

	// Register the Lambda trigger with the Flogo app
	trg := app.NewTrigger(&lambda.LambdaTrigger{}, nil)
	trg.NewFuncHandler(nil, RunActivities)

	// Return a pointer to the app
	return app
}

// RunActivities is where the magic happens. This is where you get the input from any event that might trigger
// your Lambda function in a map called evt (which is part of the inputs). The below sample,
// will simply log "Go Serverless v1.x! Your function executed successfully!" and return the same as a response.
// The trigger, in main.go, will take care of marshalling it into a proper response for the API Gateway
func RunActivities(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	// Get the items from SSM
	in := map[string]interface{}{"action": "retrieveList", "parameterName": fmt.Sprintf("%s,%s,%s", appkey, apptoken, applist), "decryptParameter": true}
	out, err := flogo.EvalActivity(&awsssm.MyActivity{}, in)
	if err != nil {
		return nil, err
	}

	key := out["result"].Value().(map[string]interface{})[appkey].(string)
	token := out["result"].Value().(map[string]interface{})[apptoken].(string)
	list := out["result"].Value().(map[string]interface{})[applist].(string)

	// Create a TrelloEvent from the incoming event
	trelloEvent := &TrelloEvent{}
	eventMap := inputs["evt"].Value().(map[string]interface{})
	err = fillStruct(eventMap["Event"].(map[string]interface{}), trelloEvent)
	if err != nil {
		return nil, err
	}

	// Run the Trello activity
	in = map[string]interface{}{"token": token, "appkey": key, "list": list, "position": "bottom", "title": trelloEvent.Title, "description": trelloEvent.Description}
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
