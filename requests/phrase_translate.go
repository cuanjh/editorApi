package requests

type PhraseTranslate struct {
	From      string `json:"from"`   // 类型
	To        string `json:"to"`     // 要翻译的目标语言
	Parent    string `json:"parent"` // uuid
	ContentTr string `json:"contentTr"`
}

type PhraseTranslateDetailRequests struct {
	From   string `json:"from"`   // 类型
	To     string `json:"to"`     // 要翻译的目标语言
	Parent string `json:"parent"` // uuid
}

type PhraseTranslateAllRequests struct {
	From    string   `json:"from"`    // 类型
	To      string   `json:"to"`      // 要翻译的目标语言
	Parents []string `json:"parents"` // uuid
}
