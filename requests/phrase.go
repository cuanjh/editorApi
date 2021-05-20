package requests

// 短语 一对多的形式
type Phrase struct {
	From     string   `json:"from"`      //类型
	Content  string   `json:"content"`   //短语单词
	Uuid     string   `json:"uuid"`      //uuid
	DictUuid []string `json:"dict_uuid"` //多个字典UUID
}

type PhraseAllRequests struct {
	From     string   `json:"from"`      // 类型
	To       string   `json:"to"`        // 要翻译的目标语言
	DictUuid []string `json:"dict_uuid"` //多个字典UUID
}
