package requests

type Catalogs struct {
	Code          string            `bson:"code" json:"code"`
	Parent_uuid   string            `bson:"parent_uuid" json:"parent_uuid"`
	Has_changed   bool              `bson:"has_changed" json:"has_changed"`
	Update_time   int64             `bson:"update_time" json:"update_time"`
	Tags          []string          `bson:"tags" json:"tags"`
	AttrTag       string            `bson:"attr_tag" json:"attr_tag"`
	List_order    int               `bson:"list_order" json:"list_order"`
	Is_show       bool              `bson:"is_show" json:"is_show"`
	OnlineState   int8              `bson:"onlineState" json:"onlineState"`
	Title         map[string]string `bson:"title" json:"title"`
	GoalTitle     string            `bson:"goalTitle" json:"goalTitle"`
	Name          string            `bson:"name" json:"name"`
	Has_del       bool              `bson:"has_del" json:"has_del"`
	Flag          []string          `bson:"flag" json:"flag"`
	Desc          map[string]string `bson:"desc" json:"desc"`
	Uuid          string            `bson:"uuid" json:"uuid"`
	Type          string            `bson:"type" json:"type"`
	Cover         []string          `bson:"cover" json:"cover"`
	Content_model string            `bson:"content_model" json:"content_model"`
}
