package model

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
)

const DB_NAME = "go-game"
const WORDS_COLLECTION = "words"

type Word struct {
	Text     string   `bson:"text,omitempty"`
	IsWord   bool     `bson:"is_word,omitempty"`
}

type Connection struct {
	session       *mgo.Session
	words         *mgo.Collection
}

func NewConnection(uri string) (*Connection) {
	session, err := mgo.Dial(uri)

	if err != nil {
		log.Panicln("MongoDB: ", err)
	}

	ensureIndex(session, WORDS_COLLECTION,
		[]string{"text"})

	session.SetMode(mgo.Monotonic, true)

	return &Connection{
		session:       session,
		words:         session.DB(DB_NAME).C(WORDS_COLLECTION),
	}
}

func (c *Connection) CloseConnection() {
	c.session.Close()
}

func (c *Connection) UpsertWord(word Word) {
	_, err := c.words.Upsert(bson.M{"text": word.Text}, word)
	if err != nil {
		log.Printf("Upsert word error: %s", err)
	}
}

func (c* Connection) GetWord(text string) (*Word) {
	var word *Word
	c.words.Find(bson.M{"text": text}).One(&word)

	if word == nil {
		return nil
	} else {
		return word
	}
}

func ensureIndex(s *mgo.Session, table string, column []string) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(DB_NAME).C(table)

	index := mgo.Index{
		Key:        column,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}