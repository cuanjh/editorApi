package responses

/**
语言信息
**/
type CourseLangs struct {
	IsShow        bool              `bson:"is_show" json:"isShow"`
	Title         map[string]string `bson:"title" json:"title"`
	Flag          []string          `bson:"flag" json:"flag"`
	Desc          map[string]string `bson:"desc" json:"desc"`
	ListOrder     int64             `bson:"list_order" json:"listOrder"`
	LanCode       string            `bson:"lan_code" json:"lanCode"`
	WordDirection string            `bson:"word_direction" json:"wordDirection"`
	IsHot         bool              `bson:"is_hot" json:"isHot"`
	HasDel        bool              `bson:"has_del" json:"hasDel"`
}
