package url

import (
	"encore.app/model"
	"testing"
)

// Run tests using `encore test`, which compiles the Encore app and then runs `go test`.
// It supports all the same flags that the `go test` command does.
// You automatically get tracing for tests in the local dev dash: http://localhost:9400
// Learn more: https://encore.dev/docs/develop/testing
//
// TestValidateNewTick - test that tick item is valid and does not hit rule
func TestValidateNewTick(t *testing.T) {
	type TestValidate struct {
		desc          string
		state         model.GameState
		tick          Tick
		expectedError bool
	}

	testData := []TestValidate{
		{desc: "happy case", state: model.GameState{Width: 10, Height: 10}, tick: Tick{VelY: 1, VelX: 0}, expectedError: false},
		{desc: "sad case: rule 180 turns", state: model.GameState{Width: 10, Height: 10, Snake: model.Snake{
			VelX: 10, VelY: 0,
		}}, tick: Tick{VelX: -10, VelY: 0}, expectedError: true},
		{desc: "sad case: rule out of bound", state: model.GameState{Width: 10, Height: 10, Snake: model.Snake{
			X: 12, Y: 2,
		}}, tick: Tick{VelX: -10, VelY: 0}, expectedError: true},
	}

	for _, testItem := range testData {
		err := ValidateNewTick(&testItem.state, testItem.tick)
		if testItem.expectedError == false && err != nil {
			t.Errorf("got %q, want %q", err, testItem.expectedError)
		}
	}
}
