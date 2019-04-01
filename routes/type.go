package routes

import (
	"mime/multipart"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	macaron "gopkg.in/macaron.v1"
	"yegoo.com/yegoo-marking-publish/utils"
)

// Context 上下文
type Context struct {
	*macaron.Context
	Cache cache.Cache
}

// Success 成功
func (c *Context) Success(data interface{}) {
	c.JSON(http.StatusOK, JSONResult{200, "", data})
}

// Error401 成功
func (c *Context) Error401(err error) {
	c.JSON(http.StatusOK, JSONResult{401, err.Error(), nil})
}

// Error4xx 成功
func (c *Context) Error4xx(err error) {
	c.JSON(http.StatusOK, JSONResult{400, err.Error(), nil})
}

// Error5xx 成功
func (c *Context) Error5xx(err error) {
	c.JSON(http.StatusOK, JSONResult{500, err.Error(), nil})
}

// NotFound 页面未找到
func (c *Context) NotFound() {
	c.JSON(http.StatusOK, JSONResult{404, "", nil})
}

// Contexter .
func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, cache cache.Cache) {
		c := &Context{Context: ctx, Cache: cache}
		ctx.Map(c)
	}
}

// JSONResult json格式返回值
type JSONResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// 小程序baseform
type XcxBaseForm struct {
	StoreId   string `form:"storeId" binding:"Required"`   // 医院id
	VersionId string `form:"versionId" binding:"Required"` // 版本id
}

// 小程序配置form
type XcxConfigForm struct {
	XcxBaseForm
	AppId string `form:"appId" binding:"Required"` // 小程序appId
	// Logo             *multipart.FileHeader `form:"logo" binding:"Required"`             // 小程序logo图片文件 后缀必须是jpg
	Logo             string `form:"logo" binding:"Required"`
	TabFont          string `form:"tabFont" binding:"Required"`          // tab文字内容(5个值 逗号分隔)
	TabColor         string `form:"tabColor" binding:"Required"`         // tab文字的默认颜色
	TabSelectedColor string `form:"tabSelectedColor" binding:"Required"` // tab文字的选中颜色
	Host             string `form:"host" binding:"Required"`             // 医院运营系统域名
}

// 公库新增版本form
type AddXcxProjectForm struct {
	VersionId   string                `form:"versionId" binding:"Required"` // 版本id
	ProjectFile *multipart.FileHeader `form:"file" binding:"Required"`      // 项目代码(压缩包)
}

// 自定义参数校验器 扩展了默认的校验规则 所有参数绑定方式(Bind BindIgnErr Form MultiPartForm Json)都可以直接使用
func (xcf XcxConfigForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {

	if errs != nil && len(errs) > 0 {
		return errs
	}

	var colorPattern = "^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$"

	match, _ := regexp.MatchString(colorPattern, xcf.TabColor)
	if !match {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"TabColor"},
			Classification: "ValidationError",
			Message:        "color value must be hexadecimal.",
		})
		return errs
	}

	match, _ = regexp.MatchString(colorPattern, xcf.TabSelectedColor)
	if !match {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"TabSelectedColor"},
			Classification: "ValidationError",
			Message:        "color value must be hexadecimal.",
		})
		return errs
	}

	arr := strings.Split(xcf.TabFont, ",")
	if len(arr) != 5 {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"TabFont"},
			Classification: "ValidationError",
			Message:        "tab character must have 5 elements which is split by ','.",
		})
		return errs
	}
	for _, ele := range arr {
		if ele == "" {
			errs = append(errs, binding.Error{
				FieldNames:     []string{"TabFont"},
				Classification: "ValidationError",
				Message:        "tab character elements must not be null.",
			})
			return errs
		}
	}

	// fileName := xcf.Logo.Filename
	// suffix := path.Ext(fileName)
	// if suffix != ".jpg" {
	// 	errs = append(errs, binding.Error{
	// 		FieldNames:     []string{"Logo"},
	// 		Classification: "ValidationError",
	// 		Message:        "only support .jpg",
	// 	})
	// 	return errs
	// }
	// size := xcf.Logo.Size
	// if size > 10240 { // 文件的最大大小为10K
	// 	errs = append(errs, binding.Error{
	// 		FieldNames:     []string{"Logo"},
	// 		Classification: "ValidationError",
	// 		Message:        "file max size is 10240 Byte, current file size is " + utils.Int64ToStr(size) + " Byte.",
	// 	})
	// 	return errs
	// }

	return nil
}

func (axpf AddXcxProjectForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	if errs != nil && len(errs) > 0 {
		return errs
	}

	fileName := axpf.ProjectFile.Filename // 上传文件名
	suffix := path.Ext(fileName)          // 上传文件后缀
	if suffix != ".zip" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"ProjectFile"},
			Classification: "ValidationError",
			Message:        "only support .zip",
		})
		return errs
	}

	return nil
}

// 重写参数校验异常处理器 只有参数绑定方式为Bind时可以使用 因为只有Bind方式提供了默认的校验异常处理 其他方式必须自行处理异常
func (xbf XcxBaseForm) Error(ctx *macaron.Context, errs binding.Errors) {
	ValidErrorHandler(ctx, errs)
}

// 自定义校验异常处理器 所有参数绑定方式(Bind BindIgnErr Form MultiPartForm Json)都可以直接使用
func ValidErrorHandler(ctx *macaron.Context, errs binding.Errors) {
	if errs == nil || len(errs) == 0 {
		return
	}
	for _, err := range errs {
		message := "[" + err.Kind() + "]: Field " + utils.Lcfirst(err.FieldNames[0]) + " is Invalid(" + err.Message + ")"
		ctx.JSON(http.StatusOK, JSONResult{400, message, nil})
		return
	}

}
