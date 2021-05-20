package editor

var TbCourseContentInfos string = "course_content_infos"

type Course_content_infos struct {
	Cover       []string            `bson:"cover" json:"cover"`
	Is_show     bool                `bson:"is_show" json:"is_show"`
	Title       map[string]string   `bson:"title" json:"title"`
	Name        string              `bson:"name" json:"name"`
	Flag        []string            `bson:"flag" json:"flag"`
	Has_changed bool                `bson:"has_changed" json:"has_changed"`
	Code        string              `bson:"code" json:"code"`
	Has_del     bool                `bson:"has_del" json:"has_del"`
	Lan_code    string              `bson:"lan_code" json:"lan_code"`
	Course_code string              `bson:"course_code" json:"course_code"`
	Desc        map[string]string   `bson:"desc" json:"desc"`
	Uuid        string              `bson:"uuid" json:"uuid"`
	Version     string              `bson:"version" json:"version"`
	Parent_uuid string              `bson:"parent_uuid" json:"parent_uuid"`
	Update_time int64               `bson:"update_time" json:"update_time"`
	Authorities []*CatalogAuthority `bson:"authorities" json:"authorities"`
	Module      string              `bson:"module" json:"module"`
}
