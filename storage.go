package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
)

var db *gorm.DB

type Chat struct {
	gorm.Model
	ChatId          int64 `sql:"index; unique"`
	UserId          int64
	ChatMin         float32
	ChatLastCommand string
	CommandDate     time.Time
}

func (c *Chat) CancelCommand() {
	c.ChatLastCommand = ""
	c.Save()
}

func (c *Chat) Save() {
	db.Save(c)
}

func initStorage() {

	var err error
	db, err = gorm.Open("sqlite3", "db.sqlite")

	if err != nil {
		log.Fatal("gorm open ", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Chat{})
}
