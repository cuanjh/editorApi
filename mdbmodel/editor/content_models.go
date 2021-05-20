package editor

type Content_models struct {
	Desc      string      `bson:"desc" json:"desc"`
	Model_key string      `bson:"model_key" json:"model_key"`
	Feilds    []FeildInfo `bson:"feilds" json:"feilds"`
	Name      string      `bson:"name" json:"name"`
}

type FeildInfo struct {
	ListOrder int         `bson:"list_order" json:"list_order"`
	Feild     string      `bson:"feild" json:"feild"`
	Type      string      `bson:"type" json:"type"`
	Name      string      `bson:"name" json:"name"`
	DataFrom  string      `bson:"data_from" json:"data_from"`
	Desc      string      `bson:"desc" json:"desc"`
	SubFeilds []FeildInfo `bson:"sub_feilds" json:"sub_feilds"`
}
