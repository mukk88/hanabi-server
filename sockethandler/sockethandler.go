package sockethandler

import (
	"log"
	"encoding/json"
	"github.com/mukk88/hanabi-server/gamelogic"
	"github.com/googollee/go-socket.io"	
	"net/http"	
)

// SocketHandler is
type SocketHandler struct {
	server *socketio.Server
} 

// NewSocketHandler is the constructor for SocketHandler
func NewSocketHandler() *SocketHandler {
	sh := SocketHandler{}
	sh.createServer()
	return &sh
} 

func (sh *SocketHandler) createServer() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		panic(err)
	}
	sh.server = server
}

func statusToResponse(status string) string {
	response := StatusResponse{status}
	responseAsString, err := json.Marshal(response)
	if err != nil {
		return ""
	}
	return string(responseAsString)
}

func createGameResponse(status string, token string) string {
	response := CreateGameResponse{status, token}
	responseAsString, err := json.Marshal(response)
	if err != nil {
		return ""
	}
	return string(responseAsString)
}

// HandleConnections handles all connections for the server
func (sh *SocketHandler) HandleConnections() {
	sh.server.On("connection", func(so socketio.Socket) {
		
		so.On("join game", func(msg string) string {
			var joinGameParams JoinGameParams
			err := json.Unmarshal([]byte(msg), &joinGameParams)
			if err != nil{
				log.Println(err)
				return statusToResponse("fail")
			}
			so.Join(joinGameParams.Token)
			return sh.joinGame(joinGameParams)
		})
		so.On("create game", func(msg string) string {
			var createGameParams CreateGameParams
			err := json.Unmarshal([]byte(msg), &createGameParams)
			if err != nil || createGameParams.Name == ""{
				return createGameResponse("fail", "")
			}
			return sh.createGame(createGameParams)
		})
		so.On("action discard", func(msg string) string {
			var discardActionParams PlayDiscardActionParams
			err := json.Unmarshal([]byte(msg), &discardActionParams)
			if err != nil {
				return statusToResponse("fail")
			}
			return sh.actionDiscard(discardActionParams)
		})
		so.On("action play", func(msg string) string {
			var playActionParams PlayDiscardActionParams
			err := json.Unmarshal([]byte(msg), &playActionParams)
			if err != nil {
				return statusToResponse("fail")
			}
			return sh.actionPlay(playActionParams)
		})
		so.On("action clue", func(msg string) string {
			var clueActionParams ClueActionParams
			err := json.Unmarshal([]byte(msg), &clueActionParams)
			if err != nil {
				return statusToResponse("fail")
			}
			return sh.actionClue(clueActionParams)
		})
		so.On("all games", func(msg string) string {
			so.Join("AllGames")
			sh.sendAllGames()
			return statusToResponse("success")
		})

		so.On("refresh game", func(msg string) string {
			var refreshGameParams RefreshGameParams
			err := json.Unmarshal([]byte(msg), &refreshGameParams)
			if err != nil {
				return statusToResponse("fail")
			}
			return sh.refreshGame(refreshGameParams)
		})

		so.On("leave game", func(msg string) {

		})

		so.On("delete game", func(msg string) {

		})

		so.On("disconnection", func() {
			log.Println("socket disconnected")
		})
	})
}

func (sh *SocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sh.server.ServeHTTP(w ,r)
}

func getAllGamesMetadata() string {
	log.Println("Getting all game metadata")
	dataStore := NewDataStore()
	defer dataStore.closeSession()
	gameMetadata, err := dataStore.allGameMetadata()
	if err != nil {
		gameMetadata = []gamelogic.Game {} 
	}
	log.Println(len(gameMetadata))
	for i:=0; i<len(gameMetadata); i++ {
		log.Println(gameMetadata[i].MetaData.Name)
	}
	value, err := json.Marshal(gameMetadata)	
	return string(value)
}

func (sh *SocketHandler) sendAllGames() {
	sh.server.BroadcastTo("AllGames", "game created", getAllGamesMetadata())
}

func (sh *SocketHandler) createGame(params CreateGameParams) string {
	log.Println("Creating game..")
	dataStore := NewDataStore()
	defer dataStore.closeSession()
	token := generateToken(7)
	newGame := gamelogic.SetupGame(params.Name, token, params.AllowedPlayers)
	err := dataStore.insertGame(&newGame)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	log.Println("Game created..")
	// select only the metadata and send to all relevant clients
	sh.server.BroadcastTo("AllGames", "game created", getAllGamesMetadata())
	return createGameResponse("success", token)
}

func (sh *SocketHandler) joinGame(params JoinGameParams) string {
	log.Println("Joining game..")
	dataStore := NewDataStore()
	defer dataStore.closeSession()
	token := params.Token
	selectedGame, err := dataStore.findGameByToken(token)
	if err != nil {
		log.Println(err)
		return statusToResponse("fail")
	}
	editedGame, isAdded := gamelogic.AddPlayer(selectedGame, params.PlayerName)
	if isAdded {
		log.Println("added")
		_, err := dataStore.updateGame(&editedGame, token)
		if err != nil {
			log.Println(err)	
			return statusToResponse("fail")
		}
	} else {
		log.Println("not added")
		return statusToResponse("fail")
	}
	log.Println("Game joined..")
	gameAsJSON, _ := json.Marshal(editedGame)
	sh.server.BroadcastTo(token, "game changed", string(gameAsJSON))
	return statusToResponse("success")
}

func (sh *SocketHandler) actionClue(params ClueActionParams) string {
	log.Println("Executing clue..")
	dataStore := NewDataStore()
	defer dataStore.closeSession()
	token := params.Token
	selectedGame, err := dataStore.findGameByToken(token)
	if err != nil {
		return statusToResponse("fail")
	}
	editedGame, _ := gamelogic.Clue(selectedGame, params.FromPlayerIndex, params.ToPlayerIndex, params.CardIndexes, params.IsNum)
	log.Println("Clue executed..")
	_, err = dataStore.updateGame(&editedGame, token)
	if err != nil {
		log.Println(err)	
		return statusToResponse("fail")
	}
	gameAsJSON, _ := json.Marshal(editedGame)
	sh.server.BroadcastTo(token, "game changed", string(gameAsJSON))
	return statusToResponse("success")
}

func (sh *SocketHandler) actionDiscard(params PlayDiscardActionParams) string {
	log.Println("Executing discard..")
	dataStore := NewDataStore()
	defer dataStore.closeSession()
	token := params.Token
	selectedGame, err := dataStore.findGameByToken(token)
	if err != nil {
		log.Println(err)
		return statusToResponse("fail")
	}
	editedGame, _ := gamelogic.Discard(selectedGame, params.FromPlayerIndex, params.CardIndex)
	log.Println("Discard executed..")
	_, err = dataStore.updateGame(&editedGame, token)
	if err != nil {
		log.Println(err)	
		return statusToResponse("fail")
	}
	gameAsJSON, _ := json.Marshal(editedGame)
	sh.server.BroadcastTo(token, "game changed", string(gameAsJSON))
	return statusToResponse("success")
}

func (sh *SocketHandler) actionPlay(params PlayDiscardActionParams) string {
	log.Println("Executing play..")
	dataStore := NewDataStore()
	defer dataStore.closeSession()
	token := params.Token
	selectedGame, err := dataStore.findGameByToken(token)
	if err != nil {
		return statusToResponse("fail")
	}
	editedGame, _ := gamelogic.Play(selectedGame, params.FromPlayerIndex, params.CardIndex)
	log.Println("Play executed..")
	_, err = dataStore.updateGame(&editedGame, token)
	if err != nil {
		log.Println(err)	
		return statusToResponse("fail")
	}
	gameAsJSON, _ := json.Marshal(editedGame)
	sh.server.BroadcastTo(token, "game changed", string(gameAsJSON))
	return statusToResponse("success")
}

func (sh *SocketHandler) refreshGame(params RefreshGameParams) string {
	log.Println("Refreshing game..")
	dataStore := NewDataStore()
	defer dataStore.closeSession()
	token := params.Token
	selectedGame, err := dataStore.findGameByToken(token)
	if err != nil {
		return statusToResponse("fail")
	}
	gameAsJSON, _ := json.Marshal(selectedGame)
	sh.server.BroadcastTo(token, "game changed", string(gameAsJSON))
	return statusToResponse("success")
}