package editor

type DiscoverChannel struct {
	Title     map[string]string `bson:"title" json:"title"`
	Icon      string            `bson:"icon" json:"icon"`
	ListOrder int               `bson:"listOrder" json:"listOrder"`
	UUID      string            `bson:"uuid" json:"uuid"`
	ShowPos   string            `bson:"showPos" json:"showPos"`
	IsShow    bool              `bson:"isShow" json:"isShow"`
	IsDel     bool              `bson:"isDel" json:"isDel"`
}
