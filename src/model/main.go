package model

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
)

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

	ensureIndex(session, "words",
		[]string{"text"})

	session.SetMode(mgo.Monotonic, true)

	return &Connection{
		session:       session,
		words:         session.DB("go-game").C("words"),
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

func (c* Connection) GetWord() (text string) {
	var word *Word
	c.words.Find(bson.M{
		"text": text, "channel_id": channelID,
	}).One(&message)

	if message == nil {
		return nil
	}

	var channel *Channel
	c.channel.Find(bson.M{
		"chat.id": channelID,
	}).One(&channel)
	message.Channel = channel

	return message
}

func ensureIndex(s *mgo.Session, table string, column []string) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("embed-telebot").C(table)

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