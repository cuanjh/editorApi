package service

import (
	"editorApi/config"
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tiw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tiw/v20190919"
)

type CourseFilesService struct {
}

var credential *common.Credential
var cpf *profile.ClientProfile

func (s *CourseFilesService) SetCredential(ctx *gin.Context) {
	// 必要步骤：
	// 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId，secretKey。
	// 你也可以直接在代码中写死密钥对，但是小心不要将代码复制、上传或者分享给他人，
	// 以免泄露密钥对危及你的财产安全。
	credential = common.NewCredential(
		config.GinVueAdminconfig.Tencent.SecretId,
		config.GinVueAdminconfig.Tencent.SecretKey,
	)
	// 非必要步骤
	// 实例化一个客户端配置对象，可以指定超时时间等配置
	cpf = profile.NewClientProfile()

	// SDK默认使用POST方法。
	// 如果你一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求。
	// 如非必要请不要修改默认设置。
	cpf.HttpProfile.ReqMethod = "POST"
	// SDK有默认的超时时间，如非必要请不要修改默认设置。
	// 如有需要请在代码中查阅以获取最新的默认值。
	cpf.HttpProfile.ReqTimeout = 30
	// SDK会自动指定域名。通常是不需要特地指定域名的，但是如果你访问的是金融区的服务，
	// 则必须手动指定域名，例如云服务器的上海金融区域名： cvm.ap-shanghai-fsi.tencentcloudapi.com
	cpf.HttpProfile.Endpoint = "tiw.tencentcloudapi.com"
	// SDK默认用HmacSHA256进行签名，它更安全但是会轻微降低性能。
	// 如非必要请不要修改默认设置。
	cpf.SignMethod = "HmacSHA1"
	// SDK 默认用 zh-CN 调用返回中文。此外还可以设置 en-US 返回全英文。
	// 但大部分产品或接口并不支持全英文的返回。
	// 如非必要请不要修改默认设置。
	cpf.Language = "en-US"
	//打印日志，默认是false
	// cpf.Debug = true
}

func (s *CourseFilesService) SetTranscodeCallback(ctx *gin.Context) (response *tiw.SetTranscodeCallbackResponse, err error) {
	s.SetCredential(ctx)
	client, _ := tiw.NewClient(credential, config.GinVueAdminconfig.Tencent.Region, cpf)
	request := tiw.NewSetTranscodeCallbackRequest()
	request.Callback = common.StringPtr(config.GinVueAdminconfig.Tencent.TranscodeFileCallbackUrl)
	request.SdkAppId = common.Int64Ptr(config.GinVueAdminconfig.Tencent.SdkAppId)
	response, err = client.SetTranscodeCallback(request)
	return
}

func (s *CourseFilesService) CreateTranscode(ctx *gin.Context, params requests.CourseFilesRequests) (response *tiw.CreateTranscodeResponse, err error) {
	s.SetCredential(ctx)
	client, _ := tiw.NewClient(credential, config.GinVueAdminconfig.Tencent.Region, cpf)
	request := tiw.NewCreateTranscodeRequest()

	request.SdkAppId = common.Int64Ptr(config.GinVueAdminconfig.Tencent.SdkAppId)
	request.Url = common.StringPtr(config.GinVueAdminconfig.UploadConfig.UploadAssets + params.FileUrl)
	
	response, err = client.CreateTranscode(request)

	if err != nil {
		return
	}

	var model = repository.CourseFiles{}
	params.TaskId = *response.Response.TaskId
	model.Create(ctx, params)
	return
}

func (s *CourseFilesService) DescribeTranscodeCallback(ctx *gin.Context) (response *tiw.DescribeTranscodeCallbackResponse, err error) {
	s.SetCredential(ctx)
	client, _ := tiw.NewClient(credential, config.GinVueAdminconfig.Tencent.Region, cpf)

	request := tiw.NewDescribeTranscodeCallbackRequest()
	request.SdkAppId = common.Int64Ptr(config.GinVueAdminconfig.Tencent.SdkAppId)

	response, err = client.DescribeTranscodeCallback(request)
	return
}

func (s *CourseFilesService) List(ctx *gin.Context, params requests.CourseFilesListRequests) (result []responses.CourseFilesListResponse, err error) {
	s.SetCredential(ctx)
	var model = repository.CourseFiles{}
	result, err = model.List(ctx, params)
	return
}

func (s *CourseFilesService) DeleteFile(ctx *gin.Context, params requests.CourseFilesDeleteRequests) (result interface{}, err error) {
	var model = repository.CourseFiles{}
	result, err = model.DeleteFile(ctx, params)
	return
}

func (s *CourseFilesService) DescribeTranscode(ctx *gin.Context, params requests.CourseFilesTranscodeRequests) (response *tiw.DescribeTranscodeResponse, err error) {
	s.SetCredential(ctx)
	client, _ := tiw.NewClient(credential, config.GinVueAdminconfig.Tencent.Region, cpf)
	request := tiw.NewDescribeTranscodeRequest()
	request.SdkAppId = common.Int64Ptr(config.GinVueAdminconfig.Tencent.SdkAppId)
	request.TaskId = common.StringPtr(params.TaskId)
	response, err = client.DescribeTranscode(request)
	if err != nil {
		return
	}

	if *response.Response.Status == "FINISHED" {
		var model = repository.CourseFiles{}
		var paramsData = requests.CourseFilesEventData{
			CompressFileURL: *response.Response.CompressFileUrl,
			ResultUrl:       *response.Response.ResultUrl,
			Pages:           *response.Response.Pages,
			Progress:        *response.Response.Pages,
			Resolution:      *response.Response.Resolution,
			TaskId:          *response.Response.TaskId,
			Title:           *response.Response.Title,
			Status:          *response.Response.Status,
		}

		_, err = model.UpdateEventData(ctx, paramsData)
	}

	return
}
