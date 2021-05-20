package responses

import "time"

// 短语 一对多的形式
type PhraseResponse struct {
	Uuid      string    `bson:"uuid" json:"uuid"`            //uuid
	Content   string    `bson:"content" json:"content"`      //短语单词
	ContentTr string    `json:"contentTr"`                   //短语翻译
	DictUuid  []string  `bson:"dict_uuid" json:"dict_uuid"`  //多个字典UUID
	IsDel     bool      `bson:"is_del" json:"isDel"`         //是否删除
	CreatedOn time.Time `bson:"created_on" json:"createdOn"` //创建时间
	UpdatedOn time.Time `bson:"updated_on" json:"updatedOn"` //更新时间
}

type PhraseTranslateResponse struct {
	Parent    string `bson:"parent" json:"parent"`        // uuid
	ContentTr string `bson:"content_tr" json:"contentTr"` //翻译内容
}
