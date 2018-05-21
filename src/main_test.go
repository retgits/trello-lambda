package main

import (
	"encoding/json"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("Successful Request", func(t *testing.T) {
		byteArray := []byte(`{"EventVersion": "1.0", "EventSource": "aws:lambda", "Trello": {"Title": "Hello World", "Description": "Hello World is a great way to test things"}}`)
		var datamap map[string]interface{}
		if err := json.Unmarshal(byteArray, &datamap); err != nil {
			panic(err)
		}

		err := handler(datamap)
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})
}
