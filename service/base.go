package service

type BaseService struct {
	TeacherService           *TeacherService           `inject:""`
	TeacherAuditService      *TeacherAuditService      `inject:""`
	CourseFilesService       *CourseFilesService       `inject:""`
	QRcodeService            *QRcodeService            `inject:""`
	DictService              *DictService              `inject:""`
	DictTranslateService     *DictTranslateService     `inject:""`
	SentenceService          *SentenceService          `inject:""`
	SentenceTranslateService *SentenceTranslateService `inject:""`
	StatisticService         *StatisticService         `inject:""`
	ReportsService           *ReportsService           `inject:""`
	ActorsService            *ActorsService            `inject:""`
	ContentReportsService    *ContentReportsService    `inject:""`
}
