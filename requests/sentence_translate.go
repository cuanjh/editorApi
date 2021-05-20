package requests

type SentenceTranslate struct {
	From      string `json:"from"`   // 类型
	To        string `json:"to"`     // 要翻译的目标语言
	Parent    string `json:"parent"` // uuid
	ContentTr string `json:"contentTr"`
	Tags      []Tag `bson:"tags" json:"tags"` // 多语言标签
}