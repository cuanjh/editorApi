package requests

type ActorsCreateRequests struct {
	Name      string  `json:"name" bson:"name"`            // 声优名称
	Photo     string  `json:"photo" bson:"photo"`          // 头像
	Sound     string  `json:"sound" bson:"sound"`          // 声音地址
	Gender    int     `json:"gender" bson:"gender"`        // 性别 1,男；0，女；
	SoundTime float64 `json:"soundTime" bson:"sound_time"` // 时长
	Role      string  `json:"role" bson:"role"`            // 分组
	Country   string  `json:"country" bson:"country"`      // 国籍
	City      string  `json:"city" bson:"city"`            // 城市
	Feature   string  `json:"feature" bson:"feature"`      // 声音特点
	Lang      string  `json:"lang" bson:"lang"`            // 语言
	Desc      string  `json:"desc" bson:"desc"`            // 描述
	Status    int     `json:"status" bson:"status"`        // 状态 1，上线；2，下线
}

type ActorsFindRequests struct {
	Uuid      string  `json:"uuid" bson:"uuid"`            // Uuid
	Name      string  `json:"name" bson:"name"`            // 声优名称
	Photo     string  `json:"photo" bson:"photo"`          // 头像
	Sound     string  `json:"sound" bson:"sound"`          // 声音地址
	Gender    int     `json:"gender" bson:"gender"`        // 性别 1,男；0，女；
	SoundTime float64 `json:"soundTime" bson:"sound_time"` // 时长
	Role      string  `json:"role" bson:"role"`            // 分组
	Country   string  `json:"country" bson:"country"`      // 国籍
	City      string  `json:"city" bson:"city"`            // 城市
	Feature   string  `json:"feature" bson:"feature"`      // 声音特点
	Lang      string  `json:"lang" bson:"lang"`            // 语言
	Desc      string  `json:"desc" bson:"desc"`            // 描述
	Status    int     `json:"status" bson:"status"`        // 状态 1，上线；2，下线
}

type ActorsFindOneRequests struct {
	Uuid      string  `json:"uuid" bson:"uuid"`            // Uuid
	Name      string  `json:"name" bson:"name"`            // 声优名称
	Photo     string  `json:"photo" bson:"photo"`          // 头像
	Sound     string  `json:"sound" bson:"sound"`          // 声音地址
	Gender    int     `json:"gender" bson:"gender"`        // 性别 1,男；0，女；
	SoundTime float64 `json:"soundTime" bson:"sound_time"` // 时长
	Role      string  `json:"role" bson:"role"`            // 分组
	Country   string  `json:"country" bson:"country"`      // 国籍
	City      string  `json:"city" bson:"city"`            // 城市
	Feature   string  `json:"feature" bson:"feature"`      // 声音特点
	Lang      string  `json:"lang" bson:"lang"`            // 语言
	Desc      string  `json:"desc" bson:"desc"`            // 描述
	Status    int     `json:"status" bson:"status"`        // 状态 1，上线；2，下线
}

type ActorsListRequests struct {
	Page
	Uuid      string  `json:"uuid" bson:"uuid"`            // Uuid
	Name      string  `json:"name" bson:"name"`            // 声优名称
	Photo     string  `json:"photo" bson:"photo"`          // 头像
	Sound     string  `json:"sound" bson:"sound"`          // 声音地址
	Gender    int     `json:"gender" bson:"gender"`        // 性别 1,男；0，女；
	SoundTime float64 `json:"soundTime" bson:"sound_time"` // 时长
	Role      string  `json:"role" bson:"role"`            // 分组
	Country   string  `json:"country" bson:"country"`      // 国籍
	City      string  `json:"city" bson:"city"`            // 城市
	Feature   string  `json:"feature" bson:"feature"`      // 声音特点
	Lang      string  `json:"lang" bson:"lang"`            // 语言
	Desc      string  `json:"desc" bson:"desc"`            // 描述
	Status    int     `json:"status" bson:"status"`        // 状态 1，上线；2，下线
}

type ActorsUpdateRequests struct {
	Uuid      string  `json:"uuid" bson:"uuid"`            // Uuid
	Name      string  `json:"name" bson:"name"`            // 声优名称
	Photo     string  `json:"photo" bson:"photo"`          // 头像
	Sound     string  `json:"sound" bson:"sound"`          // 声音地址
	Gender    int     `json:"gender" bson:"gender"`        // 性别 1,男；0，女；
	SoundTime float64 `json:"soundTime" bson:"sound_time"` // 时长
	Role      string  `json:"role" bson:"role"`            // 分组
	Country   string  `json:"country" bson:"country"`      // 国籍
	City      string  `json:"city" bson:"city"`            // 城市
	Feature   string  `json:"feature" bson:"feature"`      // 声音特点
	Lang      string  `json:"lang" bson:"lang"`            // 语言
	Desc      string  `json:"desc" bson:"desc"`            // 描述
	Status    int     `json:"status" bson:"status"`        // 状态 1，上线；2，下线
}

type ActorsDeleteRequests struct {
	Uuid      string  `json:"uuid" bson:"uuid"`            // Uuid
	Name      string  `json:"name" bson:"name"`            // 声优名称
	Photo     string  `json:"photo" bson:"photo"`          // 头像
	Sound     string  `json:"sound" bson:"sound"`          // 声音地址
	Gender    int     `json:"gender" bson:"gender"`        // 性别 1,男；0，女；
	SoundTime float64 `json:"soundTime" bson:"sound_time"` // 时长
	Role      string  `json:"role" bson:"role"`            // 分组
	Country   string  `json:"country" bson:"country"`      // 国籍
	City      string  `json:"city" bson:"city"`            // 城市
	Feature   string  `json:"feature" bson:"feature"`      // 声音特点
	Lang      string  `json:"lang" bson:"lang"`            // 语言
	Desc      string  `json:"desc" bson:"desc"`            // 描述
	Status    int     `json:"status" bson:"status"`        // 状态 1，上线；2，下线
}
