package wx

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bluecatlee/WechatDevToolsHttpCall/utils"
)

const (
	ide            = "C://Users//Administrator//AppData//Local//微信web开发者工具//User Data//Default//.ide"                                 //端口文件
	openUrl        = "http://127.0.0.1:{port}/open?projectpath={projectpath}"                                                         // 打开项目请求路径
	loginUrl       = "http://127.0.0.1:{port}/login?format=base64"                                                                    // 获取登录二维码请求路径
	previewUrl     = "http://127.0.0.1:{port}/preview?projectpath={projectpath}&format=base64"                                        // 获取预览二维码请求路径
	uploadUrl      = "http://127.0.0.1:{port}/upload?projectpath={projectpath}&version={version}&desc={desc}&infooutput={infooutput}" // 上传请求路径
	closeUrl       = "http://127.0.0.1:{port}/close?projectpath={projectpath}"                                                        // 关闭当前项目窗口
	quitUrl        = "http://127.0.0.1:{port}/quit"                                                                                   // 关闭开发者工具
	buildnpmUrl    = "http://127.0.0.1:{port}/buildnpm?projectpath={projectpath}&compiletype=miniprogram"                             // 构建npm
	testUrl        = "http://127.0.0.1:{port}/test?projectpath={projectpath}"                                                         // 自动化测试
	autopreviewUrl = "http://127.0.0.1:{port}/autopreview?projectpath={projectpath}"                                                  // 自动预览
)

// 打开小程序项目 (打开不同的项目会打开多个窗口)
func Open(projectpath string) (err error) {
	if len([]rune(strings.Trim(projectpath, " "))) == 0 {
		return fmt.Errorf("projectpath must not be empty")
	}

	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return err
	}

	// 参数进行编码
	projectpath = utils.UrlEncode(projectpath)

	// 发起请求
	url := strings.Replace(openUrl, "{port}", pid, -1)
	url = strings.Replace(url, "{projectpath}", projectpath, -1)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(string(data))
	}

	return nil
}

// 获取登录二维码 (调用后所有窗口的登录状态都会失效,重新登录成功后所有窗口都会登录这个微信账号)
func Login() (qrCode string, err error) {
	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return "", err
	}

	// 发起请求
	url := strings.Replace(loginUrl, "{port}", pid, -1)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(string(data))
	}

	return string(data), nil

}

// 获取预览二维码
func Preview(projectpath, previewPage string) (qrCode string, err error) {
	if len([]rune(strings.Trim(projectpath, " "))) == 0 {
		return "", fmt.Errorf("projectpath must not be empty")
	}

	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return "", err
	}

	projectpath = utils.UrlEncode(projectpath)

	// 发起请求
	url := strings.Replace(previewUrl, "{port}", pid, -1)
	url = strings.Replace(url, "{projectpath}", projectpath, -1)
	if previewPage != "" {
		url = url + `&compilecondition={"pathName":"` + previewPage + `","query":""}` // 设置默认打开的页面
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(string(data))
	}

	qrCode = string(data)
	if !strings.HasPrefix(qrCode, "data:image/jpeg;base64,") {
		qrCode = "data:image/jpeg;base64," + qrCode
	}

	return qrCode, nil
}

// 上传小程序
func Upload(projectpath, version, desc, infooutput string) (err error) {
	if len([]rune(strings.Trim(projectpath, " "))) == 0 {
		return fmt.Errorf("projectpath must not be empty")
	}
	if len([]rune(strings.Trim(version, " "))) == 0 {
		return fmt.Errorf("version must not be empty")
	}

	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return err
	}

	projectpath = utils.UrlEncode(projectpath)

	// 发起请求
	url := strings.Replace(uploadUrl, "{port}", pid, -1)
	url = strings.Replace(url, "{projectpath}", projectpath, -1)
	url = strings.Replace(url, "{version}", version, -1)
	url = strings.Replace(url, "{desc}", desc, -1)
	url = strings.Replace(url, "{infooutput}", infooutput, -1)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(string(data))
	}

	return nil

}

// 关闭当前项目窗口
func Close(projectpath string) (err error) {
	if len([]rune(strings.Trim(projectpath, " "))) == 0 {
		return fmt.Errorf("projectpath must not be empty")
	}
	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return err
	}

	projectpath = utils.UrlEncode(projectpath)

	// 发起请求
	url := strings.Replace(closeUrl, "{port}", pid, -1)
	url = strings.Replace(url, "{projectpath}", projectpath, -1)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(string(data))
	}

	return nil
}

// 关闭开发者工具
func Quit() (err error) {
	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return err
	}
	// 发起请求
	url := strings.Replace(quitUrl, "{port}", pid, -1)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(string(data))
	}

	return nil
}

// 构建npm
func Buildnpm(projectpath string) (err error) {
	if len([]rune(strings.Trim(projectpath, " "))) == 0 {
		return fmt.Errorf("projectpath must not be empty")
	}
	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return err
	}

	projectpath = utils.UrlEncode(projectpath)

	// 发起请求
	url := strings.Replace(buildnpmUrl, "{port}", pid, -1)
	url = strings.Replace(url, "{projectpath}", projectpath, -1)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(string(data))
	}

	return nil
}

// 自动化测试
func AutoTest(projectpath string) (result string, err error) {
	if len([]rune(strings.Trim(projectpath, " "))) == 0 {
		return "", fmt.Errorf("projectpath must not be empty")
	}
	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return "", err
	}

	projectpath = utils.UrlEncode(projectpath)

	// 发起请求
	url := strings.Replace(testUrl, "{port}", pid, -1)
	url = strings.Replace(url, "{projectpath}", projectpath, -1)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(string(data))
	}

	return string(data), nil

}

// 自动预览
func AutoPreview(projectpath string) (err error) {
	if len([]rune(strings.Trim(projectpath, " "))) == 0 {
		return fmt.Errorf("projectpath must not be empty")
	}
	// 获取端口号
	pid, err := getPID()
	if err != nil {
		return err
	}

	projectpath = utils.UrlEncode(projectpath)

	// 发起请求
	url := strings.Replace(autopreviewUrl, "{port}", pid, -1)
	url = strings.Replace(url, "{projectpath}", projectpath, -1)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(string(data))
	}

	return nil
}

// 获取微信开发者工具端口号
func getPID() (pid string, err error) {
	data, err := ioutil.ReadFile(ide)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

/*
 * 异常返回示例
 *	{\"code\":40000,\"error\":\"错误 需要重新登录\"}
 *	{\"code\":40000,\"error\":\"错误 Error: 登录用户不是该小程序的开发者\"}
 *	{\"code\":40000,\"error\":\"错误 Error: 代码包大小为 2239 kb，上限为 2048 kb，请删除文件后重试\"}
 *	{\"code\":40000,\"error\":\"错误 Error: AppID 不合法\"}
 *  {\"code\":400002,\"error\":\"project.config.json 中缺少了 appid\"}
 *  {\"code\":400002,\"error\":\"没有找到 node_modules 目录。 查看文档 https://developers.weixin.qq.com/miniprogram/dev/devtools/npm.html\"}
 *  {\"code\":40000,\"error\":\"错误 Error: 今天已经提交过测试\"}
 *  ...
 */
