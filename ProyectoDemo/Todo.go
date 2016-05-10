package main

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Todo struct {
	Id        bson.ObjectId  `json:"_id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Completed bool      `json:"completed" bson:"completed"`
	Due       time.Time `json:"due" bson:"due"`
}

type Todos []Todo
