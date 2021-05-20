package editor

type Course_infos struct {
	Is_show      bool              `bson:"is_show" json:"is_show"`
	Uuid         string            `bson:"uuid" json:"uuid"`
	Name         string            `bson:"name" json:"name"`
	Code         string            `bson:"code" json:"code"`
	Flag         []string          `bson:"flag" json:"flag"`
	Course_type  int64             `bson:"course_type" json:"course_type"`
	Cover        []string          `bson:"cover" json:"cover"`
	Lan_code     string            `bson:"lan_code" json:"lan_code"`
	Desc         map[string]string `bson:"desc" json:"desc"`
	Title        map[string]string `bson:"title" json:"title"`
	HasDict      bool              `bson:"has_dict" json:"has_dict"`
	SoundActors  []soundActor      `bson:"sound_actors" json:"sound_actors"`
	DefaultActor string            `bson:"default_actor" json:"default_actor"`
}

type soundActor struct {
	Role   string `bson:"role" json:"role"`
	Name   string `bson:"name" json:"name"`
	Photo  string `bson:"photo" json:"photo"`
	Gender int    `bson:"gender" json:"gender"`
	Sound  string `bson:"sound" json:"sound"`
}
