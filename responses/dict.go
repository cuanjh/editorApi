package responses

import "time"

type DictResponse struct {
	Uuid          string           `bson:"uuid" json:"uuid"`              // uuid
	Content       string           `bson:"content" json:"content"`        // 单词
	Images        []Image          `bson:"images" json:"images"`          // 图片
	SoundInfos    []SoundInfos     `bson:"sound_infos" json:"soundInfos"` // 音标
	DictTranslate DictTranslate    `json:"dictTranslate"`                 // 翻译
	Phrase        []PhraseResponse `json:"phrase"`                        // 短语
	IsDel         bool             `bson:"is_del" json:"isDel"`           // 是否删除
	Done          bool             `bson:"done" json:"done"`              // 是否上线
	CreatedOn     time.Time        `bson:"created_on" json:"createdOn"`   // 创建时间
	UpdatedOn     time.Time        `bson:"updated_on" json:"updatedOn"`   // 更新时间
}

type Image struct {
	Url  string `bson:"url" json:"url"`
	Name string `bson:"name" json:"name"`
}

type SoundInfos struct {
	Ct     string `bson:"ct" json:"ct"`         // 类型 en：英； ：美;
	Ps     string `bson:"ps" json:"ps"`         // 音标
	Sound  string `bson:"sound" json:"sound"`   // 声音
	Gender string `bson:"gender" json:"gender"` // male: 男音 female：女音
}

type DictTranslate struct {
	Parent    string      `bson:"parent" json:"parent"`        // uuid
	Expansion string      `bson:"expansion" json:"expansion"`  // 拓展
	ContentTr []ContentTr `bson:"content_tr" json:"contentTr"` // 词义
	Synonym   []Synonym   `bson:"synonym" json:"synonym"`      // 近义词
	Homonyms  []Homonyms  `bson:"homonyms" json:"homonyms"`    // 同词根
	Tags      []Tag       `bson:"tags" json:"tags"`            // 多语言标签
}

type ContentTr struct {
	Cx      string `json:"cx"`      // 词性
	Content string `json:"content"` // 词义
}

type WordAttr struct {
	Content   string `json:"content"`   // 单词
	ContentTr string `json:"contentTr"` // 词义
}

type Homonyms struct {
	Cx    string     `json:"cx"`    //词性
	Attrs []WordAttr `json:"attrs"` //词义
}

type Synonym struct {
	Cx        string `json:"cx"`        //词性
	Content   string `json:"content"`   //单词
	ContentTr string `json:"contentTr"` //词义
}
