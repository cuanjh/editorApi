package requests

/**
语言信息
**/
type CourseLangs struct {
	IsShow        bool              `json:"is_show"`
	Title         map[string]string `json:"title"`
	Flag          []string          `json:"flag"`
	Desc          map[string]string `json:"desc"`
	ListOrder     int64             `json:"list_order"`
	LanCode       string            `json:"lan_code"`
	WordDirection string            `json:"word_direction"`
	IsHot         bool              `json:"is_hot"`
	HasDel        bool              `json:"has_del"`
}
