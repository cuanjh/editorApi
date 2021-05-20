package editor

type Course_model_video struct {
	Parent_uuid string   `bson:"parent_uuid" json:"parent_uuid"`
	Tags        []string `bson:"tags" json:"tags"`
	Video       string   `bson:"video" json:"video"`
	Has_del     bool     `bson:"has_del" json:"has_del"`
	List_order  int      `bson:"list_order" json:"list_order"`
	Has_changed bool     `bson:"has_changed" json:"has_changed"`
	Cover       string   `bson:"cover" json:"cover"`
	Content     string   `bson:"content" json:"content"`
	Update_time int      `bson:"update_time" json:"update_time"`
	Video_time  float64  `bson:"video_time" json:"video_time"`
	Uuid        string   `bson:"uuid" json:"uuid"`
}
