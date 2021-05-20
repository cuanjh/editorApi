package requests

type CommentsUploadHanderParams struct {
	FilePath   string `json:"file_path"`
	CourseUuid string `json:"course_uuid"` //直播ChatroomID
}
