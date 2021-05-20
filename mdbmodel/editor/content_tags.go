package editor

import "gopkg.in/mgo.v2/bson"

type Content_tags struct {
	Key        string   `bson:"key" json:"key"`
	Type       string   `bson:"type" json:"type"`
	Name       string   `bson:"name" json:"name"`
	Desc       bson.M   `bson:"desc" json:"desc"`
	Title      bson.M   `bson:"title" json:"title"`
	ListOrder  int      `bson:"list_order" json:"list_order"`
	Cover      []string `bson:"cover" json:"cover"`
	Flag       []string `bson:"flag" json:"flag"`
	HasDel     bool     `bson:"has_del" json:"has_del"`
	HasChanged bool     `bson:"has_changed" json:"has_changed"`
}
