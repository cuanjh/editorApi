package editor

const TbCatalogs = "catalogs"

type ExaminStateInfo struct {
	UserUUID string `bson:"user_id" json:"user_id"`
	Comment  string `bson:"comment" json:"comment"`
}

type CatalogAuthority struct {
	UUID            string          `bson:"uuid" json:"uuid"`
	UserUUID        string          `bson:"user_uuid" json:"user_uuid"`
	Auth            string          `bson:"auth" json:"auth"`
	ExaminState     int             `bson:"examinState" json:"examin_state"` //权限状态，0正在编辑，1提交审核，2审核通过，3审核没通过
	ExaminStateInfo ExaminStateInfo `bson:"examinStateInfo" json:"examin_state_info"`
}

type Catalogs struct {
	Parent_uuid   string             `bson:"parent_uuid" json:"parent_uuid"`
	Has_changed   bool               `bson:"has_changed" json:"has_changed"`
	Update_time   int64              `bson:"update_time" json:"update_time"`
	Tags          []string           `bson:"tags" json:"tags"`
	AttrTag       string             `bson:"attr_tag" json:"attr_tag"`
	List_order    int                `bson:"list_order" json:"list_order"`
	Is_show       bool               `bson:"is_show" json:"is_show"`
	OnlineState   int8               `bson:"onlineState" json:"onlineState"`
	Title         map[string]string  `bson:"title" json:"title"`
	GoalTitle     string             `bson:"goalTitle" json:"goalTitle"`
	Name          string             `bson:"name" json:"name"`
	Has_del       bool               `bson:"has_del" json:"has_del"`
	Flag          []string           `bson:"flag" json:"flag"`
	Desc          map[string]string  `bson:"desc" json:"desc"`
	Uuid          string             `bson:"uuid" json:"uuid"`
	Type          string             `bson:"type" json:"type"`
	Cover         []string           `bson:"cover" json:"cover"`
	Content_model string             `bson:"content_model" json:"content_model"`
	Authorities   []CatalogAuthority `bson:"authorities" json:"authorities"`
}
