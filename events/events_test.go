package events

import (
	"bytes"
	"encoding/json"
	"testing"
)

func Test_newErrorEvent_with_invalid_code(t *testing.T) {
	event := NewError(0, "The error message")

	data, _ := json.Marshal(event)
	expected := `{"event":"pusher:error","data":{"code":null,"message":"The error message"}}`

	if bytes.Compare(data, []byte(expected)) != 0 {
		t.Errorf("%s != %s", string(data), expected)
	}
}

func Test_newErrorEvent_with_valid_code(t *testing.T) {
	event := NewError(4007, "Unsupported protocol version")

	data, _ := json.Marshal(event)
	expected := `{"event":"pusher:error","data":{"code":4007,"message":"Unsupported protocol version"}}`

	if bytes.Compare(data, []byte(expected)) != 0 {
		t.Errorf("%s != %s", string(data), expected)
	}
}
