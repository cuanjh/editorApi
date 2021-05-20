package requests

import (
	"encoding/json"
	"time"
)

type DictUploadParams struct {
	Filename string `json:"filename"`
}

type DictHanderParams struct {
	FilePath string `json:"file_path"`
	From     string `json:"from"` //来自 en
	To       string `json:"to"`   //去向	 zh
}

// 字典
type Dict struct {
	From       string       `json:"from"` //类型
	CardId     string       `json:"card_id"`
	Uuid       string       `json:"uuid"`       //uuid
	ListOrder  string       `json:"listOrder"`  //排序
	Images     []Image      `json:"images"`     //图片
	Content    string       `json:"content"`    //单词
	SoundInfos []SoundInfos `json:"soundInfos"` //音标
	IsDel      bool         `json:"is_del"`     //是否删除
	CreatedOn  time.Time    `json:"created_on"` //创建时间
}

// 音标
type SoundInfos struct {
	Ct     string `json:"ct"`     //类型 en：英；us：美;
	Ps     string `json:"ps"`     //音标
	Sound  string `json:"sound"`  //声音
	Gender string `json:"gender"` //male: 男音 female：女音
}

//上线下线
type DictOnlineRequests struct {
	From  string   `json:"from" validate:"required" label:"原单词"`   //来自 en
	To    string   `json:"to" validate:"required" label:"目标语言"`    //去向	 zh
	CType string   `json:"ctype" validate:"required" label:"内容类型"` //word or sentence
	Uuids []string `json:"uuids" validate:"required" label:"内容的UUID的数组"`
}

type DictListRequests struct {
	Page
	SearchType int    `json:"searchType"`                           //搜索类型：0（精确），1（模糊）
	Content    string `json:"content"`                              //搜索内容
	From       string `json:"from" validate:"required" label:"原单词"` //来自 en
	To         string `json:"to" validate:"required" label:"目标语言"`  //去向	 zh
	IsDel      bool   `json:"is_del"`                               //是否删除
	OnLine     string `json:"online"`                               //是否删除
}

type DictDetailRequests struct {
	Uuid string `json:"uuid" validate:"required" label:"uuid"` //uuid
	From string `json:"from" validate:"required" label:"原单词"`  //来自 en
	To   string `json:"to" validate:"required" label:"目标语言"`   //去向	 zh
}

type DictUpdateRequests struct {
	From          string           `json:"from" validate:"required" label:"原单词"` //来自 en
	To            string           `json:"to" validate:"required" label:"目标语言"`  //去向 zh
	Uuid          string           `bson:"uuid" json:"uuid"`                     // uuid
	Content       string           `bson:"content" json:"content"`               // 单词
	Images        []Image          `bson:"images" json:"images"`                 // 图片
	SoundInfos    []SoundInfos     `bson:"sound_infos" json:"soundInfos"`        // 音标
	DictTranslate DictTranslate    `json:"dictTranslate"`                        // 翻译
	Phrase        []PhraseResponse `json:"phrase"`                               // 短语
	IsDel         bool             `bson:"is_del" json:"isDel"`                  // 是否删除
}
type DictDelRequests struct {
	From  string   `json:"from" validate:"required" label:"原单词"` //来自 en
	To    string   `json:"to" validate:"required" label:"目标语言"`  //去向 zh
	CType string   `json:"ctype"`
	Uuids []string `bson:"uuids" json:"uuids"` // uuid
}

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

type Image struct {
	Url  string `bson:"url" json:"url"`
	Name string `bson:"name" json:"name"`
}

type DictSound struct {
	Uuid       string       `json:"uuid" validate:"required" label:"uuid"` //uuid
	From       string       `json:"from" validate:"required" label:"原单词"`  //来自 en
	SoundInfos []SoundInfos `json:"soundInfos"`                            //音标
}

type DictCardId struct {
	Uuid       string       `json:"uuid" validate:"required" label:"uuid"`      //uuid
	From       string       `json:"from" validate:"required" label:"原单词"`       //来自 en
	CardId     string       `json:"card_id" validate:"required" label:"CardId"` //来自 CartId
	Content    string       `json:"content"`                                    //内容
	SoundInfos []SoundInfos `json:"soundInfos"`                                 //音标
}

type DictAddTag struct {
	From     string `json:"from" validate:"required" label:"原单词"` //来自 en
	To       string `json:"to" validate:"required" label:"目标语言"`  //去向 zh
	FilePath string `json:"file_path"`                            //url
	Tag      string `json:"tag"`                                  //url
}

type DictFindAll struct {
	Page
	From      string    `json:"from" validate:"required" label:"原单词"` //来自 en
	CreatedOn time.Time `json:"createdOn"`                            //创建时间
}

/****************************************************** ES操作 ******************************************************/

type DictESParams struct {
	CType   string   `json:"ctype"`   //word or sentence
	Uuids   []string `json:"uuids"`   //uuid
	From    string   `json:"from"`    //来自 en
	To      string   `json:"to"`      //去向 zh
	Operate string   `json:"operate"` //操作方式 delete offline online
}

type DictSearch struct {
	ID         string          `json:"id"`
	IsWord     bool            `json:"is_word"`
	IsSentence bool            `json:"is_sentence"`
	Content    string          `json:"content"`
	CradID     string          `json:"card_id"`
	CType      string          `json:"ctype"`
	Images     []string        `json:"images"`
	Tags       []Tag           `json:"tags"`
	LangCode   string          `json:"lang_code"`
	Extras     json.RawMessage `json:"extras"`
}

type searchWordExtras struct {
	SoundInfos []*wordSound      `json:"sound_infos"`
	ContentTrs []*searchWordTr   `json:"content_trs"`
	Synonyms   []*relevanceWord  `json:"synonyms"`
	Homonyms   []*searchRootWord `json:"homonyms"`
	Phrases    []*phraseInfo     `json:"phrases"`
}

type searchSentenceExtras struct {
	Source      string   `json:"source"`
	CourseCodes []string `json:"from_course_codes"`
	SoundInfos  []*struct {
		Gender string `json:"gender"`
		Sound  string `json:"sound"`
	} `json:"sound_infos"`
	ContentTr string `json:"content_tr"`
}

type wordSound struct {
	Ps     string `json:"ps"` //音标
	Ct     string `json:"ct"` //国家,en/us
	Sound  string `json:"sound"`
	Gender string `json:"gender"`
	Photo  string `json:"photo"`
}

type searchWordTr struct {
	Cx string `json:"cx"`
	Tr string `json:"content"`
}

type relevanceWord struct {
	Cx        string `json:"cx"`
	Content   string `json:"content"`
	ContentTr string `json:"contentTr"`
}

type searchRootWord struct {
	Cx        string          `json:"cx"`
	RootWords []*rootWordInfo `json:"attrs"`
}

type rootWordInfo struct {
	Content   string `json:"content"`
	ContentTr string `json:"contentTr"`
}

type phraseInfo struct {
	Phrase   string `json:"phrase"`
	PhraseTr string `json:"phraseTr"`
}
