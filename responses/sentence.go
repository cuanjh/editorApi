package responses

import "time"

type Sentence struct {
	Uuid        string                `bson:"uuid" json:"uuid"`                //uuid
	Mold        int                   `bson:"mold" json:"mold"`                //类型 1：口语；2：书面语；
	Sentence    string                `bson:"sentence" json:"sentence"`        //句子
	Image       []string              `bson:"image" json:"image"`              //图片
	SoundInfos  []SentenceSoundInfos  `bson:"sound_infos" json:"soundInfos"`   //声音
	CourseInfos []SentenceCourseInfos `bson:"course_infos" json:"courseInfos"` //课程
	Source      string                `bson:"source" json:"source"`            //来源
	IsDel       bool                  `bson:"is_del" json:"isDel"`             //是否删除
	SentenceTr  SentenceTranslate     `json:"sentenceTr"`                      //翻译内容
	CreatedOn   time.Time             `bson:"created_on" json:"createdOn"`     // 创建时间
	UpdatedOn   time.Time             `bson:"updated_on" json:"updatedOn"`     // 更新时间
	Done        bool                  `bson:"done" json:"done"`                // 是否上线，true上线,false下线
}

type SentenceSoundInfos struct {
	Sound  string `json:"sound"`  //声音
	Gender string `json:"gender"` //male: 男音 female：女音
}

// 短语 一对多的形式
type SentenceCourseInfos struct {
	Uuid        string `bson:"uuid" json:"uuid"`                //uuid
	CourseCode  string `bson:"course_code" json:"courseCode"`   //courseCode
	ChapterCode string `bson:"chapter_code" json:"chapterCode"` //chapterCode
	Image       string `bson:"image" json:"image"`              //image
	Sound       string `bson:"sound" json:"sound"`              //sound
}

type SentenceTranslate struct {
	Parent    string `bson:"parent" json:"parent"` // uuid
	ContentTr string `bson:"content_tr" json:"contentTr"`
	Tags      []Tag  `bson:"tags" json:"tags"` // 多语言标签
}

type Tag struct {
	Key  string `bson:"key" json:"key"`
	Name string `bson:"name"  json:"name"`
}

type SentenceFindAll struct {
	Uuid      string    `bson:"uuid" json:"uuid"`            //uuid
	Sentence  string    `bson:"sentence" json:"sentence"`    //句子
	CreatedOn time.Time `bson:"created_on" json:"createdOn"` // 创建时间
	UpdatedOn time.Time `bson:"updated_on" json:"updatedOn"` // 更新时间
}
