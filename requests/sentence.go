package requests

import "time"

type SentenceSearchRequests struct {
	Page
	From       string `json:"from" validate:"required" label:"from"` //from
	To         string `json:"to" validate:"required" label:"to"`     //to
	Sentence   string `json:"sentence"`                              //句子
	SearchType int    `json:"searchType"`                            //0 精确；1 模糊
}

// 短语 一对多的形式
type Sentence struct {
	// Page
	CardId      string               `json:"card_id"`
	From        string               `json:"from" validate:"required" label:"from"` //from
	To          string               `json:"to" validate:"required" label:"to"`     //to
	Uuid        string               `bson:"uuid" json:"uuid"`                      //uuid
	Mold        int                  `json:"mold"`                                  //类型 1：口语；2：书面语；
	Sentence    string               `bson:"sentence" json:"sentence"`              //句子
	Image       []string             `json:"image" bson:"images"`                   //图片
	SoundInfos  []SentenceSoundInfos `json:"sound_infos" json:"soundInfos"`         //声音
	CourseInfos []CourseInfos        `bson:"course_infos" json:"courseInfos"`       //课程
	Source      string               `json:"source"`                                //来源
	IsDel       bool                 `bson:"is_del" json:"isDel"`                   //是否删除
	Done        bool                 `bson:"done" json:"done"`
	CreatedOn   time.Time            `bson:"created_on" json:"createdOn"` //是否删除
	UpdatedOn   time.Time            `bson:"updated_on" json:"updatedOn"` //是否删除
}

// 短语 一对多的形式
type SentenceDetail struct {
	Page
	From        string               `json:"from" validate:"required" label:"from"` //from
	To          string               `json:"to" validate:"required" label:"to"`     //to
	Uuid        string               `json:"uuid" validate:"required" label:"uuid"` //uuid
	Mold        int                  `json:"mold"`                                  //类型 1：口语；2：书面语；
	Sentence    string               `json:"sentence"`                              //句子
	Image       []string             `json:"image"`                                 //图片
	SoundInfos  []SentenceSoundInfos `json:"soundInfos"`                            //声音
	CourseInfos []CourseInfos        `json:"courseInfos"`                           //课程
	Source      string               `json:"source"`                                //来源
	IsDel       bool                 `json:"isDel"`                                 //是否删除
	SentenceTr  SentenceTranslate    `json:"sentenceTr"`                            //翻译内容
}

type SentenceSoundInfos struct {
	Sound  string `json:"sound"`  //声音
	Gender string `json:"gender"` //male: 男音 female：女音
}

type CourseInfos struct {
	Uuid        string `json:"uuid"` //uuid
	CourseCode  string `json:"courseCode"`
	ChapterCode string `json:"chapterCode"`
	Image       string `json:"image"`
	Sound       string `json:"sound"`
}

type SentenceFindAll struct {
	CreatedOn time.Time `json:"createdOn"`                             //创建时间
	From      string    `json:"from" validate:"required" label:"from"` //from
	Page
}

type SentenceUpdate struct {
	From        string               `json:"from" validate:"required" label:"from"`   //from
	To          string               `json:"to" validate:"required" label:"to"`       //to
	Uuid        string               `json:"uuid" validate:"required" label:"uuid"`   //uuid
	Mold        int                  `json:"mold"`                                    //类型 1：口语；2：书面语；
	Sentence    string               `json:"sentence" validate:"required" label:"句子"` //句子
	Image       []string             `json:"image"`                                   //图片
	SoundInfos  []SentenceSoundInfos `json:"soundInfos"`                              //声音
	CourseInfos []CourseInfos        `json:"courseInfos"`                             //课程
	Source      string               `json:"source"`                                  //来源
	IsDel       bool                 `json:"isDel"`                                   //是否删除
	SentenceTr  SentenceTranslate    `json:"sentenceTr"`                              //翻译内容
}

type SentenceDelete struct {
	From string `json:"from" validate:"required" label:"from"` //from
	Uuid string `json:"uuid" validate:"required" label:"uuid"` //uuid
}

type SentenceCardId struct {
	From   string `json:"from" validate:"required" label:"from"` //from
	Uuid   string `json:"uuid" validate:"required" label:"uuid"` //uuid
	CardId string `json:"card_id" validate:"required" label:"card_id"`
}
