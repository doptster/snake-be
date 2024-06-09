package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encore.dev/storage/sqldb"
	mathRand "math/rand"
	"time"

	"encore.app/constant"
	"encore.app/model"
)

// NewGameParams is the request data for the new endpoint.
type NewGameParams struct {
	Width  int `query:"w"`
	Height int `query:"h"`
}

// generateGameID generates a random short game ID.
func generateGameID() (string, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

// generateFruitPosition generates a random fruit position on the board
func generateFruitPosition(width, height int) (int, int) {
	// Create a new random source and generator
	source := mathRand.NewSource(time.Now().UnixNano())
	random := mathRand.New(source)

	// Generate random x and y coordinates
	x := random.Intn(width)
	y := random.Intn(height)

	return x, y
}

// insertGame inserts a new game into the database.
func insertGame(ctx context.Context, state *model.GameState) error {
	_, err := db.Exec(ctx, `
        INSERT INTO game (id, width, height, score, fruitX, fruitY)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, state.GameID, state.Width, state.Height, state.Score, state.Fruit.X, state.Fruit.Y)
	return err
}

// resetState resets the game state to new
func resetState(state *model.GameState, width int, height int) {
	fruitX, fruitY := generateFruitPosition(width, height)
	state.Fruit = model.Fruit{
		X: fruitX,
		Y: fruitY,
	}
	state.Snake = model.Snake{
		X:    constant.SnakeStartX,
		Y:    constant.SnakeStartY,
		VelX: constant.SnakeStartVelX,
		VelY: constant.SnakeStartVelY,
	}
}

// New generates a new game for the board.
//
//encore:api public method=GET path=/new
func New(ctx context.Context, params *NewGameParams) (*model.GameState, error) {
	id, err := generateGameID()
	if err != nil {
		return nil, err
	}

	state := &model.GameState{
		GameID: id,
		Width:  params.Width,
		Height: params.Height,
		Score:  constant.GameScoreDefault,
	}
	resetState(state, params.Width, params.Height)

	if err := insertGame(ctx, state); err != nil {
		return nil, err
	}

	return state, nil
}

// Define a database named 'url', using the database
// migrations  in the "./migrations" folder.
// Encore provisions, migrates, and connects to the database.
// Learn more: https://encore.dev/docs/primitives/databases
var db = sqldb.NewDatabase("url", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
