package editor

const TbCourseFiles = "course_files"

type CourseFiles struct {
	Uuid                 string                `bson:"uuid" json:"uuid"`
	LiveUuid             string                `bson:"live_uuid" json:"live_uuid"`
	FileUrl              string                `bson:"file_url" json:"file_url"`
	Type                 string                `bson:"type" json:"type"`
	Size                 int64                 `bson:"size" json:"size"`
	CreatedOn            string                `bson:"created_on" json:"created_on"`
	State                int                   `bson:"state" json:"state"`
	Title                string                `bson:"title" json:"title"`
	TaskId               string                `bson:"task_id" json:"task_id"`
	CourseFilesEventData *CourseFilesEventData `bson:"event_data" json:"event_data"`
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