package sockethandler

// ClueActionParams is
type ClueActionParams struct {
	Token string
	FromPlayerIndex int
	ToPlayerIndex int
	CardIndexes []int
	IsNum bool
}

// PlayDiscardActionParams is
type PlayDiscardActionParams struct {
	Token string
	FromPlayerIndex int
	CardIndex int
}

// JoinGameParams is
type JoinGameParams struct {
	Token string
	PlayerName string
}

// CreateGameParams is 
type CreateGameParams struct {
	Name string
	AllowedPlayers int
}

// RefreshGameParams is 
type RefreshGameParams struct {
	Token string
}

// StatusResponse is
type StatusResponse struct {
	Status string
}

// CreateGameResponse is 
type CreateGameResponse struct {
	Status string
	Token string
}

