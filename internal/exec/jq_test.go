package exec_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sven-seyfert/apiprobe/internal/exec"
)

const inputJSON = `
[
    {
        "id": 101,
        "name": "Alice",
        "address": {
            "street": "Hauptstr. 5",
            "city": "Berlin",
            "zip": "10115"
        },
        "projects": [
            {
                "title": "Website Redesign",
                "status": "in progress"
            },
            {
                "title": "Mobile App",
                "status": "completed"
            },
            {
                "title": "Native App",
                "status": "done"
            }
        ]
    }
]
`

func TestGoJQ_simpleQuery(t *testing.T) {
	jqCommand := ".[0].address.city"
	expected := `"Berlin"`

	received, err := exec.GoJQ(context.Background(), jqCommand, []byte(inputJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertJSONEqual(t, []byte(expected), received)
}

func TestGoJQ_containsQuery(t *testing.T) {
	jqCommand := ".[0].projects[] | select(.title | contains(\"App\")) | .title"
	expected := `["Mobile App","Native App"]`

	received, err := exec.GoJQ(context.Background(), jqCommand, []byte(inputJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertJSONEqual(t, []byte(expected), received)
}

func TestGoJQ_emptyJqArgs(t *testing.T) {
	jqCommand := ""
	expected := inputJSON

	received, err := exec.GoJQ(context.Background(), jqCommand, []byte(inputJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertJSONEqual(t, []byte(expected), received)
}

func TestGoJQ_invalidJSONInput(t *testing.T) {
	jqCommand := ".[0].address.city"
	invalidJSONInput := `
[
    {
        "id": 101,
        "address": {
            "street": "Hauptstr. 5",
            "city": "Berlin"
            "zip": "10115"
        }
    }
]
`

	_, err := exec.GoJQ(context.Background(), jqCommand, []byte(invalidJSONInput))
	if err == nil {
		t.Fatalf("no error but one is expected for this scenario")
	}
}

func TestGoJQ_invalidJQCommand(t *testing.T) {
	invalidJQCommand := ".[0].address city"

	_, err := exec.GoJQ(context.Background(), invalidJQCommand, []byte(inputJSON))
	if err == nil {
		t.Fatalf("no error but one is expected for this scenario")
	}
}

func assertJSONEqual(t *testing.T, expected, received []byte) {
	t.Helper()

	var expectedVal, receivedVal interface{}

	if err := json.Unmarshal(expected, &expectedVal); err != nil {
		t.Fatalf("invalid expected JSON: %v", err)
	}

	if err := json.Unmarshal(received, &receivedVal); err != nil {
		t.Fatalf("invalid received JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedVal, receivedVal) {
		expPretty, _ := json.MarshalIndent(expectedVal, "", "  ")
		recPretty, _ := json.MarshalIndent(receivedVal, "", "  ")

		t.Errorf("JSON not equal:\nexpected:\n%s\n\nreceived:\n%s",
			expPretty, recPretty)
	}
}
