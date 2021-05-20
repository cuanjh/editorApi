package editor

type Course_model_kid_sound struct {
	Update_time  int     `bson:"update_time" json:"update_time"`
	Sentence     string  `bson:"sentence" json:"sentence"`
	Sound_time   float64 `bson:"sound_time" json:"sound_time"`
	Parent_uuid  string  `bson:"parent_uuid" json:"parent_uuid"`
	Has_changed  bool    `bson:"has_changed" json:"has_changed"`
	List_order   int     `bson:"list_order" json:"list_order"`
	Sentence_trs string  `bson:"sentence_trs" json:"sentence_trs"`
	Uuid         string  `bson:"uuid" json:"uuid"`
	Has_del      bool    `bson:"has_del" json:"has_del"`
	Sound        string  `bson:"sound" json:"sound"`
	Image        string  `bson:"image" json:"image"`
}
