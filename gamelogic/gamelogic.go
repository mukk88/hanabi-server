package gamelogic

import (
	"math/rand"
)

func createNewDeck() []CardData {
	deck := []CardData {}
	shuffledDeck := []CardData {}
	numSuits := 5
	numOnes := 3
	numOthers := 2
	for i := 0; i < numSuits; i++ {
		for j := 0; j < numOnes; j++ { deck = append(deck, CardData{i, 1, false, false}) }
		for j := 0; j < numOthers; j++ { deck = append(deck, CardData{i, 2, false, false}, CardData{i, 3, false, false}, CardData{i, 4, false, false}) }
		deck = append(deck, CardData{i, 5, false, false})
	}
	perm := rand.Perm(len(deck))
	for _, value := range perm {
		shuffledDeck = append(shuffledDeck, deck[value])
	}
	return shuffledDeck
}

// SetupGame performs all the setup needed for a new game
func SetupGame(name string, token string, players int) Game {
	return Game{ 
		token,
		GameMetaData{
			name,
			players,
			[]string {},
			false,
		}, 
		GameData{
			[]PlayerData {},
			0,
			createNewDeck(),
			8,
			3,
			[]int {0, 0, 0, 0, 0},
			[]CardData {},
			0,
			"",
			0,
		},
	}
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// AddPlayer adds a player to the game and deals a hand
func AddPlayer(game Game, name string) (Game, bool) {
	editedGame := Game(game)
	if stringInSlice(name, game.MetaData.PlayerNames) {
		return editedGame, true
	}
	if len(game.MetaData.PlayerNames) >= game.MetaData.AllowedPlayers {
		return editedGame, false
	}
	editedGame.MetaData.PlayerNames = append(editedGame.MetaData.PlayerNames, name)
	var handSize int
	if editedGame.MetaData.AllowedPlayers > 3 {
		handSize = 4
	} else {
		handSize = 5
	}
	newPlayer := PlayerData{name, name, len(editedGame.Data.Players), 4, editedGame.Data.Deck[0:handSize]}
	editedGame.Data.Players = append(editedGame.Data.Players, newPlayer)
	editedGame.Data.Deck = editedGame.Data.Deck[handSize:]
	if len(game.MetaData.PlayerNames) == game.MetaData.AllowedPlayers {
		editedGame.Data.Status = Playing
	}

	return editedGame, true
}

func replaceCard(editedGame *Game, playerIndex int, cardIndex int) {
	editedGame.Data.Players[playerIndex].Cards = append(
		editedGame.Data.Players[playerIndex].Cards[:cardIndex],
		editedGame.Data.Players[playerIndex].Cards[cardIndex+1:]...,
	)

	if len(editedGame.Data.Deck) > 0 {
		editedGame.Data.Players[playerIndex].Cards = append(
			editedGame.Data.Players[playerIndex].Cards,
			editedGame.Data.Deck[0],
		)
		editedGame.Data.Deck = editedGame.Data.Deck[1:]
	} else {
		editedGame.Data.LastTurnCount++
		if editedGame.Data.LastTurnCount == len(editedGame.Data.Players) {
			editedGame.Data.Status = Won
			editedGame.MetaData.Complete = true
		}
	}
}

// Discard is
func Discard(game Game, fromPlayerIndex int, cardIndex int) (Game, bool) {
	editedGame := Game(game)
	fromPlayerName := editedGame.Data.Players[fromPlayerIndex].Name
	if editedGame.Data.Clues != 8 {
		editedGame.Data.Clues++
	} 
	discardedCard := editedGame.Data.Players[fromPlayerIndex].Cards[cardIndex]
	editedGame.Data.Discards = append(editedGame.Data.Discards, discardedCard)
	replaceCard(&editedGame, fromPlayerIndex, cardIndex)
	editedGame.Data.Turn = (editedGame.Data.Turn + 1) % len(editedGame.Data.Players)
	editedGame.Data.LastMove = fromPlayerName + " discarded a card"
	
	return editedGame, true
}

// Clue is
func Clue(game Game, fromPlayerIndex int, toPlayerIndex int,
	cardIndexes []int, isNum bool) (Game, bool) {
	
	editedGame := Game(game)
	fromPlayerName := editedGame.Data.Players[fromPlayerIndex].Name
	toPlayerName := editedGame.Data.Players[toPlayerIndex].Name
	for i := 0; i < len(cardIndexes); i++ {
		if (isNum) {
			editedGame.Data.Players[toPlayerIndex].Cards[cardIndexes[i]].NumberRevealed = true
		} else {
			editedGame.Data.Players[toPlayerIndex].Cards[cardIndexes[i]].SuitRevealed = true			
		}
	}
	editedGame.Data.LastMove = fromPlayerName + " gave " + toPlayerName + " a clue"
	editedGame.Data.Clues--
	editedGame.Data.Turn = (editedGame.Data.Turn + 1) % len(editedGame.Data.Players) 
	if len(editedGame.Data.Deck) == 0 && editedGame.Data.LastTurnCount != 0 {
		editedGame.Data.LastTurnCount++
		if editedGame.Data.LastTurnCount == len(editedGame.Data.Players) {
			editedGame.Data.Status = Won
			editedGame.MetaData.Complete = true
		}
	}
	return editedGame, true
}

// Play is
func Play(game Game, fromPlayerIndex int, cardIndex int) (Game, bool) {
	editedGame := Game(game)

	valid := false
	fromPlayerName := editedGame.Data.Players[fromPlayerIndex].Name	
	playedCard := editedGame.Data.Players[fromPlayerIndex].Cards[cardIndex] 

	for i := 0; i < len(editedGame.Data.Table); i++ {
		if i == playedCard.Suit && editedGame.Data.Table[i] == playedCard.Number - 1 {
			valid = true
			editedGame.Data.Table[i]++
			if playedCard.Number == 5 && editedGame.Data.Clues != 8 {
				editedGame.Data.Clues++
			}
			break
		}
	}

	if !valid {
		editedGame.Data.Discards = append(editedGame.Data.Discards, playedCard)
		editedGame.Data.Burns--
		if editedGame.Data.Burns == 0 {
			editedGame.Data.Status = Lost
			editedGame.MetaData.Complete = true
		}
	}
	replaceCard(&editedGame, fromPlayerIndex, cardIndex)
	editedGame.Data.Turn = (editedGame.Data.Turn + 1) % len(editedGame.Data.Players)
	editedGame.Data.LastMove = fromPlayerName + " played a card"

	return editedGame, true
}