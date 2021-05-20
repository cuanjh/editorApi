package editor

type Content_model_pro_sound struct {
	Has_changed          bool     `bson:"has_changed" json:"has_changed"`
	Uuid                 string   `bson:"uuid" json:"uuid"`
	Parent_uuid          string   `bson:"parent_uuid" json:"parent_uuid"`
	List_order           int      `bson:"list_order" json:"list_order"`
	Code                 string   `bson:"code" json:"code"`
	Sentence_temp        string   `bson:"sentence_temp" json:"sentence_temp"`
	Options              []string `bson:"options" json:"options"`
	Sound                string   `bson:"sound" json:"sound"`
	Type                 string   `bson:"type" json:"type"`
	Options_phoneticize  []string `bson:"options_phoneticize" json:"options_phoneticize"`
	Sentence_phoneticize string   `bson:"sentence_phoneticize" json:"sentence_phoneticize"`
	Image                string   `bson:"image" json:"image"`
	Sentence             string   `bson:"sentence" json:"sentence"`
	Has_del              bool     `bson:"has_del" json:"has_del"`
}
