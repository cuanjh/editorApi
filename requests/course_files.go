package requests

import "time"

type CourseFilesTranscodeCallbackRequests struct {
	EventData struct {
		CompressFileURL string `json:"CompressFileUrl"`
		Error           struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
		Pages               int64  `json:"Pages"`
		Resolution          string `json:"Resolution"`
		ResultURL           string `json:"ResultUrl"`
		TaskID              string `json:"TaskId"`
		ThumbnailResolution string `json:"ThumbnailResolution"`
		ThumbnailURL        string `json:"ThumbnailUrl"`
		Title               string `json:"Title"`
	} `json:"EventData"`
	EventType  string `json:"EventType"`
	ExpireTime int64  `json:"ExpireTime"`
	SdkAppID   int64  `json:"SdkAppId"`
	Sign       string `json:"Sign"`
	Timestamp  int64  `json:"Timestamp"`
}

type CourseFilesTranscodeRequests struct {
	TaskId string `json:"task_id"`
}

type CourseFilesCreateTranscodeRequests struct {
	FileUrl  string `json:"file_url"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	Title    string `json:"title"`
	LiveUuid string `json:"live_uuid"`
}

type CourseFilesRequests struct {
	Uuid                 string               `bson:"uuid" json:"uuid"`
	LiveUuid             string               `bson:"live_uuid" json:"live_uuid"`
	FileUrl              string               `bson:"file_url" json:"file_url"`
	Type                 string               `bson:"type" json:"type"`
	Size                 int64                `bson:"size" json:"size"`
	CreatedOn            time.Time            `bson:"created_on" json:"created_on"`
	State                int                  `bson:"state" json:"state"`
	Title                string               `bson:"title" json:"title"`
	TaskId               string               `bson:"task_id" json:"task_id"`
	CourseFilesEventData CourseFilesEventData `bson:"event_data" json:"event_data"`
}

type CourseFilesEventData struct {
	CompressFileURL string `bson:"compress_file_url" json:"compress_file_url"`
	ResultUrl       string `bson:"result_url" json:"result_url"`
	Pages           int64  `bson:"pages" json:"pages"`
	Progress        int64  `bson:"progress" json:"progress"`
	Resolution      string `bson:"resolution" json:"resolution"`
	TaskId          string `bson:"task_id" json:"task_id"`
	Title           string `bson:"title" json:"title"`
	Status          string `bson:"status" json:"status"`
}

type CourseFilesListRequests struct {
	LiveUuid string `bson:"live_uuid" json:"live_uuid"`
}

type CourseFilesDeleteRequests struct {
	Uuid     string `bson:"uuid" json:"uuid" validate:"required" label:"文件UUID"`
	LiveUuid string `bson:"live_uuid" json:"live_uuid" validate:"required" label:"课程UUID"`
}
