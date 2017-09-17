package gamelogic

// CardData is data about a card
type CardData struct {
	Suit int
	Number int
	SuitRevealed bool
	NumberRevealed bool
}

// PlayerData is data about a player
type PlayerData struct {
	Name string
	Token string
	Index int
	HandSize int
	Cards []CardData
}

// GameMetaData is a struct that holds metadata
type GameMetaData struct {
	Name string
	AllowedPlayers int
	PlayerNames []string
	Complete bool
}

// PlayState is the different game states the game can be in
type PlayState int
const (
	// WaitingForPlayers is	
	WaitingForPlayers PlayState = iota
	// ReadyToPlay is
	ReadyToPlay
	// Playing is
	Playing
	// Won is
	Won
	// Lost is
	Lost
)

// GameData is a struct that holds data
type GameData struct {
	Players []PlayerData
	Turn int
	Deck []CardData
	Clues int
	Burns int
	Table []int
	Discards []CardData
	Status PlayState
	LastMove string
	LastTurnCount int
}

// Game is a struct that holds all data about the game
type Game struct {
	Token string
	MetaData GameMetaData
	Data GameData
}