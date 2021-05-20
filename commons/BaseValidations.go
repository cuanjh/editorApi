package commons

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

type BaseValidations struct {
}

func init() {

}

func (v *BaseValidations) Check(model interface{}) (message string, err error) {
	//中文翻译器
	zh_ch := zh.New()
	uni := ut.New(zh_ch)
	trans, _ := uni.GetTranslator("zh")
	//验证器
	validate := validator.New()

	//注册一个函数，获取struct tag里自定义的label作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name:=fld.Tag.Get("label")
		return name
	})

	//验证器注册翻译器
	zh_translations.RegisterDefaultTranslations(validate, trans)

	err = validate.Struct(model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			//翻译错误信息
			message = err.Translate(trans)
		}
	}
	return
}
