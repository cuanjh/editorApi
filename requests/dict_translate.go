package requests

type DictTranslateDetailRequests struct {
	From   string `json:"from"`                 // 类型
	To     string `json:"to"`                   // 要翻译的目标语言
	Parent string `bson:"parent" json:"parent"` // uuid
}

type DictTranslate struct {
	From      string      `json:"from"`                        // 类型
	To        string      `json:"to"`                          // 要翻译的目标语言
	Parent    string      `bson:"parent" json:"parent"`        // uuid
	Expansion string      `bson:"expansion" json:"expansion"`  // 拓展
	ContentTr []ContentTr `bson:"content_tr" json:"contentTr"` // 词义
	Synonym   []Synonym   `bson:"synonym" json:"synonym"`      // 近义词
	Homonyms  []Homonyms  `bson:"homonyms" json:"homonyms"`    // 同词根
	Tags      []Tag       `bson:"tags" json:"tags"`            // 多语言标签
}

type AddTags struct {
	From   string `json:"from"`                 // 类型
	Parent string `json:"parent"` // uuid
	To     string `json:"to"`                   // 要翻译的目标语言
	Tags   []Tag  `json:"tags"`     // 多语言标签
}

type Tag struct {
	Key  string `bson:"key" json:"key"`
	Name string `bson:"name"  json:"name"`
}

type ContentTr struct {
	Cx      string `json:"cx"`
	Content string `json:"content"`
}

type WordAttr struct {
	Cx        string `json:"cx"`
	Content   string `json:"content"`
	ContentTr string `json:"contentTr"`
}

type Homonyms struct {
	Cx    string     `json:"cx"`
	Attrs []WordAttr `json:"attrs"`
}

type Synonym struct {
	Cx        string `json:"cx"`
	Content   string `json:"content"`
	ContentTr string `json:"contentTr"`
}

// 词典翻译
type DictInterpretation struct {
	Translate string `json:"translate"` //词义
	Character string `json:"character"` //词性
}
