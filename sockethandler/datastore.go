package sockethandler

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/mukk88/hanabi-server/gamelogic"	
	"math/rand"
	"time"
)

// DataStore is the database wrapper
type DataStore struct {
	session *mgo.Session
} 

func (ds *DataStore) setupSession() {
	session, err := mgo.Dial("localhost:27017")
	// session, err := mgo.Dial("mongodb://hanabi:hanabi123@ds133044.mlab.com:33044/hanabi")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	ds.session = session
}

func (ds *DataStore) insertGame(newGame *gamelogic.Game) error {
	collection := ds.session.DB("hanabi").C("games")
	return collection.Insert(newGame)
}

func (ds *DataStore) findGameByToken(token string) (gamelogic.Game, error) {
	result := gamelogic.Game{}
	collection := ds.session.DB("hanabi").C("games")
	err := collection.Find(bson.M{"token": token}).One(&result)
	return result, err
}

func(ds *DataStore) updateGame(editedGame *gamelogic.Game, token string) (*mgo.ChangeInfo, error) {
	collection := ds.session.DB("hanabi").C("games")
	return collection.Upsert(
		bson.M{"token": token},
		editedGame,
	)
}

func(ds *DataStore) allGameMetadata() ([]gamelogic.Game, error) {
	result := []gamelogic.Game{}
	collection := ds.session.DB("hanabi").C("games")
	err := collection.Find(bson.M{"metadata.complete": false}).Select(bson.M{"metadata": 1, "token": 1}).All(&result)
	return result, err	
}

func (ds *DataStore) closeSession() {
	ds.session.Close()
}

// NewDataStore is the constructor for DataStore
func NewDataStore() *DataStore {
	ds := DataStore{}
	ds.setupSession()
	return &ds
} 

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateToken(n int) string {
	rand.Seed(time.Now().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}