package url

import (
	"context"
	"encore.app/constant"
	"encore.app/model"
	"encore.dev/beta/errs"
)

// Tick is the request data for the tick data.
type Tick struct {
	VelX int `json:"velX,omitempty"`
	VelY int `json:"velY,omitempty"`
}

// BatchUpdateParams is the request data for the validate endpoint.
type BatchUpdateParams struct {
	Ticks  []Tick `json:"ticks"`
	GameID string `json:"gameId"`
}

// ValidateNewTick inserts a validate a tick item
func ValidateNewTick(state *model.GameState, tick Tick) error {
	// Rule: Diagonal
	if tick.VelX != 0 && tick.VelY != 0 {
		return &errs.Error{
			Code:    errs.FailedPrecondition,
			Message: constant.GameOverMessage,
		}
	}
	// Rule: 180 turns
	if (tick.VelX > 0 && state.Snake.VelX < 0) ||
		(tick.VelX < 0 && state.Snake.VelX > 0) ||
		(tick.VelY < 0 && state.Snake.VelY > 0) ||
		(tick.VelY > 0 && state.Snake.VelY < 0) {
		return &errs.Error{
			Code:    errs.FailedPrecondition,
			Message: constant.GameOverMessage,
		}
	}
	// Rule: Not moving
	if tick.VelX == 0 && tick.VelY == 0 {
		return &errs.Error{
			Code:    errs.FailedPrecondition,
			Message: constant.GameOverMessage,
		}
	}
	newX := state.Snake.X + tick.VelX
	newY := state.Snake.Y + tick.VelY
	// Rule: Out of bound
	if newX < 0 || newX >= state.Width || newY < 0 || newY >= state.Height {
		return &errs.Error{
			Code:    errs.FailedPrecondition,
			Message: constant.GameOverMessage,
		}
	}
	// Valid tick, update virtual state
	state.Snake = model.Snake{
		VelX: tick.VelX,
		VelY: tick.VelY,
		X:    newY,
		Y:    newY,
	}
	return nil
}

// updateGame updates game data upon validate success
func updateGame(ctx context.Context, state *model.GameState) error {
	_, err := db.Exec(ctx, `
        UPDATE game SET score = $2, fruitX = $3, fruitY = $4 WHERE id = $1
    `, state.GameID, state.Score, state.Fruit.X, state.Fruit.Y)
	return err
}

// Validate checks for ticks and send a new fruit position
//
//encore:api public method=POST path=/validate
func Validate(ctx context.Context, params *BatchUpdateParams) (*model.GameState, error) {
	state := &model.GameState{
		Snake: model.Snake{
			VelY: constant.SnakeStartVelX,
			VelX: constant.SnakeStartVelY,
		},
	}
	// Fetch the latest state of game
	err := db.QueryRow(ctx, `
        SELECT id, width, height, score, fruitX, fruitY FROM game
        WHERE id = $1
    `, params.GameID).Scan(&state.GameID, &state.Width, &state.Height, &state.Score, &state.Fruit.X, &state.Fruit.Y)
	if err != nil {
		return nil, err
	}

	// Ensure each tick item is valid
	isValidTick := true
	for _, tick := range params.Ticks {
		err = ValidateNewTick(state, tick)
		if err != nil {
			return nil, err
		}
	}
	// Rule: Snake must reach Fruit position
	if state.Snake.X != state.Fruit.X || state.Snake.Y != state.Fruit.Y {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "Fruit not found, the ticks do not lead the snake to the fruit position.",
		}
	}

	// Update game state with valid ticks
	if isValidTick {
		resetState(state, state.Width, state.Height)
		state.Score++
		err = updateGame(ctx, state)
		if err != nil {
			return nil, err
		}
	}

	return state, err
}
