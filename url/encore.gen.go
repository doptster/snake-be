// Code generated by encore. DO NOT EDIT.

package url

import (
	"context"
	model "encore.app/model"
)

// These functions are automatically generated and maintained by Encore
// to simplify calling them from other services, as they were implemented as methods.
// They are automatically updated by Encore whenever your API endpoints change.

// Interface defines the service's API surface area, primarily for mocking purposes.
//
// Raw endpoints are currently excluded from this interface, as Encore does not yet
// support service-to-service API calls to raw endpoints.
type Interface interface {
	// New generates a new game for the board.
	New(ctx context.Context, p *NewGameParams) (*model.GameState, error)

	// Validate checks for ticks and send a new fruit position
	Validate(ctx context.Context, p *BatchUpdateParams) (*model.GameState, error)
}
