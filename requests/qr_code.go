package requests

type QRcodeAddRequests struct {
	Title string `json:"title"` //标题
	Info  string `json:"info"`  //json内容
	Size  int    `json:"size"`  //二维码图片大小
}

type QRcodeUpdateRequests struct {
	UUID  string `json:"uuid"  validate:"required" label:"UUID"` //UUID
	Title string `json:"title"`                                  //标题
	Info  string `json:"info"`                                   //json内容
	Size  int    `json:"size"`                                   //二维码图片大小
}

type QRcodeDeleteRequests struct {
	UUID string `json:"uuid"  validate:"required" label:"UUID"` //UUID
}

type QRcodeDetailsRequests struct {
	UUID           string `json:"uuid"  validate:"required" label:"UUID"` //UUID
	UserAgent      string `json:"userAgent"`                              //UserAgent
	AddressContent string `json:"addressContent"`                         //AddressContent
}

type QRcodeListRequests struct {
	Page
}

type QRcodeImageRequests struct {
	Size int    `json:"size"`                                //二维码图片大小
	Url  string `json:"url" validate:"required" label:"Url"` //Url内容
}
