package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type list struct {
	Content string
	Done    bool
}

type Note struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Email  string             `bson:"email,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Text   string             `bson:"text,omitempty"`
	Lists  []list             `bson:"list,omitempty"`
	Images []string           `bson:"images,omitempty"`
}
