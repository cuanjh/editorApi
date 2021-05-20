package editor

import (
	bson "go.mongodb.org/mongo-driver/bson/primitive"
)

type Course_langs struct {
	Is_show        bool              `bson:"is_show" json:"is_show"`
	Title          map[string]string `bson:"title" json:"title"`
	Flag           []string          `bson:"flag" json:"flag"`
	Desc           map[string]string `bson:"desc" json:"desc"`
	List_order     int64             `bson:"list_order" json:"list_order"`
	Lan_code       string            `bson:"lan_code" json:"lan_code"`
	Id             bson.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Word_direction string            `bson:"word_direction" json:"word_direction"`
	Is_hot         bool              `bson:"is_hot" json:"is_hot"`
	Has_del        bool              `bson:"has_del" json:"has_del"`
}
