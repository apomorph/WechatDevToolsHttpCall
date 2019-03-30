package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"yegoo.com/yegoo-marking-publish/conf"
	"yegoo.com/yegoo-marking-publish/utils"
	"yegoo.com/yegoo-marking-publish/wx"
)

const (
	indexPath    = "/pages/index/index.json"       // 首页标题配置文件路径
	goodsPath    = "/pages/goods/goods.json"       // 商品标题配置文件路径
	activityPath = "/pages/yearend/yearend.json"   // 年终庆典标题配置文件路径
	casePath     = "/pages/caselist/caselist.json" // 案例标题配置文件路径
	userPath     = "/pages/user/user.json"         // 个人中心标题配置文件路径

	logoPath = "/images/newslist/loginlogo.jpg" // logo图片

	confPath     = "/conf"     // 医院配置子目录
	confFileName = "conf.json" // 医院配置文件名字

	previewIndexPage = "pages/index/index" // 预览时打开的页面 首页
)

var titleConfigPathArr = [5]string{indexPath, goodsPath, activityPath, casePath, userPath}

// app.json文件的数据结构
type App struct {
	Pages  []string `json:"pages"`
	Window struct {
		BackgroundTextStyle          string `json:"backgroundTextStyle"`
		NavigationBarBackgroundColor string `json:"navigationBarBackgroundColor"`
		NavigationBarTitleText       string `json:"navigationBarTitleText"`
		NavigationBarTextStyle       string `json:"navigationBarTextStyle"`
	} `json:"window"`
	TabBar struct {
		List []struct {
			PagePath         string `json:"pagePath"`
			Text             string `json:"text"`
			IconPath         string `json:"iconPath"`
			SelectedIconPath string `json:"selectedIconPath"`
		} `json:"list"`
		Color         string `json:"color"`
		SelectedColor string `json:"selectedColor"`
	} `json:"tabBar"`
	NetworkTimeout struct {
		Request int `json:"request"`
	} `json:"networkTimeout"`
}

// 页面标题数据结构
type Title struct {
	Component              bool     `json:"component"`
	UsingComponents        struct{} `json:"usingComponents"`
	NavigationBarTitleText string   `json:"navigationBarTitleText"`
	EnablePullDownRefresh  bool     `json:"enablePullDownRefresh"`
}

// 自定义服务访问限制器
func AccessLimiter(c *Context) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	storeId := c.Query("storeId")
	if storeId == "" {
		err := fmt.Errorf("storeId不能为空")
		c.Error4xx(err)
		return
	}

	sessionName := conf.C.Session.Name
	sessionTimeout := conf.C.Session.Timeout

	cache := c.Cache
	val := cache.Get(sessionName)
	var sessionVal = ""
	if val != nil {
		sessionVal = val.(string) //类型强转
	}
	if sessionVal != "" && sessionVal != storeId {
		// 无法判断是否是真的登录 只能从访问级别进行限制
		err := fmt.Errorf("无法访问,服务占用中")
		// fmt.Println("storeId:" + storeId + ", 当前session:" + sessionVal + ", 无法访问")
		c.Error401(err)
		return
	} else {
		// fmt.Println("storeId:" + storeId + ", 当前session:" + sessionVal + ", 访问成功")
		cache.Put(sessionName, storeId, sessionTimeout)
	}

	return
}

// 获取登录二维码
func Login(c *Context) {
	qrCode, err := wx.Login()
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}
	c.Success(qrCode)
}

// 保存参数信息 替换小程序代码 并打开项目
func OpenProject(c *Context, form XcxConfigForm) {

	versionId := form.VersionId
	storeId := form.StoreId

	rootPath := conf.C.Respository.Root   // 小程序代码公库根目录
	storePath := conf.C.Respository.Store // 小程序代码私库根目录

	srcPath := rootPath + "/" + versionId + "/" + conf.C.Project.Name                   // 小程序代码具体版本公库路径
	destPath := storePath + "/" + storeId + "/" + versionId + "/" + conf.C.Project.Name // 具体医院的小程序代码路径

	xcxConfDir := storePath + "/" + storeId + "/" + versionId + confPath // 医院配置信息保存目录

	// First. 保存自定义配置信息
	// 参数转json
	data, err := json.Marshal(form)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}
	// 创建配置信息保存目录
	err = createConfFile(xcxConfDir, confFileName, data)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}

	// Second. 获取小程序代码并拷贝到工作目录
	// 获取对应版本的小程序代码
	err = copyXcxProject(srcPath, destPath)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}

	// Third. 修改新项目文件的自定义配置信息
	err = modifyProject(destPath, form)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}

	// Fourth. 打开微信开发者工具 并打开对应的项目
	err = wx.Open(destPath)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}

	c.Success("open project success!")

}

// 预览
func Preview(c *Context, form XcxBaseForm) {

	storeId := form.StoreId
	versionId := form.VersionId

	// 获取要预览的小程序项目所在路径
	xcxDir, err := findXcxPath(storeId, versionId)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}

	// 生成预览二维码
	qrcode, err := wx.Preview(xcxDir, previewIndexPage)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}
	c.Success(qrcode)
}

// 上传
func Upload(c *Context, form XcxBaseForm) {

	storeId := form.StoreId
	versionId := form.VersionId

	// 获取要上传的小程序项目所在路径
	xcxDir, err := findXcxPath(storeId, versionId)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}

	err = wx.Upload(xcxDir, versionId, "", "")
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}
	// 关闭当前项目窗口 TODO 无法关闭 原因未知
	// defer wx.Close(xcxDir)

	// 删除缓存中保存的storeId
	sessionName := conf.C.Session.Name
	c.Cache.Delete(sessionName)

	c.Success("上传成功")
}

// 根据医院id和版本id获取项目所在目录
func findXcxPath(storeId, versionId string) (string, error) {
	storePath := conf.C.Respository.Store
	xcxDir := storePath + "/" + storeId + "/" + versionId + "/" + conf.C.Project.Name

	// 判断路径是否存在
	b, err := utils.PathExists(xcxDir)
	if err != nil {
		return "", err
	}
	if !b {
		return "", fmt.Errorf("project not exist!")
	}
	return xcxDir, nil
}

// 创建医院配置信息文件
func createConfFile(path, filename string, data []byte) error {
	// 判断目录是否存在
	b, err := utils.PathExists(path)
	if err != nil {
		return err
	}
	// 不存在则创建
	if !b {
		err = os.MkdirAll(path, 0666)
		if err != nil {
			return err
		}
	}
	// 创建配置信息文件
	file, err := os.Create(path + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")

	return nil
}

// 小程序代码复制
func copyXcxProject(srcPath, destPath string) error {
	b, err := utils.PathExists(srcPath)
	if err != nil {
		return err
	}
	b, err = utils.PathExists(destPath)
	if err != nil {
		return err
	}
	if !b {
		err = os.MkdirAll(destPath, 0777)
		if err != nil {
			return err
		}
	}
	// 复制项目文件到对应医院的目录下
	err = utils.CopyDir(srcPath, destPath)
	if err != nil {
		return err
	}

	return nil
}

// 修改项目中的配置
func modifyProject(dest string, form XcxConfigForm) error {

	// 1)./project.config.json的appid属性值
	filePath1 := dest + "/" + "project.config.json"
	err := modifyAppId(filePath1, form.AppId)
	if err != nil {
		return err
	}

	// 2)./api/env.js的url和storeId属性值
	filePath2 := dest + "/api/" + "env.js"
	err = modifyHostAndStoreId(filePath2, form.Host, form.StoreId)
	if err != nil {
		return err
	}

	// 3)./app.json的tabBar相关属性值
	filePath3 := dest + "/" + "app.json"
	arr := strings.Split(form.TabFont, ",")
	err = modifyTabBar(filePath3, form.TabColor, form.TabSelectedColor, arr)
	if err != nil {
		return err
	}

	// 4). 部分页面标题的相关属性值修改
	for index, filePath4 := range titleConfigPathArr {
		err = modifyTitle(dest+filePath4, arr[index])
		if err != nil {
			return err
		}
	}

	// 5). 修改logo图片 todo
	filePath5 := dest + logoPath
	uploadFileReader, err := form.Logo.Open()
	if err != nil {
		return err
	}
	fileData, err := ioutil.ReadAll(uploadFileReader)
	if err != nil {
		return err
	}
	err = modifyLogo(filePath5, fileData)
	if err != nil {
		return err
	}

	return nil
}

// 修改appid
func modifyAppId(filePath, appId string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var confJson map[string]interface{}
	err = json.Unmarshal(file, &confJson)
	if err != nil {
		return err
	}
	oldAppId, ok := confJson["appid"]
	if !ok {
		err = fmt.Errorf("appid配置不存在!")
		return err
	}

	if oldAppId != appId {
		confJson["appid"] = appId
		byteValue, err := json.Marshal(confJson)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filePath, byteValue, 0777)
		if err != nil {
			return err
		}

	}

	return nil
}

// 修改host和storeId
func modifyHostAndStoreId(filePath, host, storeId string) error {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	var envStr = `let config = {
		url: 'https://` + host + `/api',
		storeId: ` + storeId + `, 
		debug: false
	  } 
	  module.exports = config`
	_, err = file.WriteString(envStr)
	if err != nil {
		return err
	}

	return nil
}

// 修改tabBar
func modifyTabBar(filePath, tabColor, selectedColor string, characters []string) error {

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var ap App
	err = json.Unmarshal(file, &ap)
	if err != nil {
		return err
	}

	ap.TabBar.Color = tabColor
	ap.TabBar.SelectedColor = selectedColor
	ap.TabBar.List[0].Text = characters[0]
	ap.TabBar.List[1].Text = characters[1]
	ap.TabBar.List[2].Text = characters[2]
	ap.TabBar.List[3].Text = characters[3]
	ap.TabBar.List[4].Text = characters[4]

	byteValue, err := json.Marshal(ap)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, byteValue, 0777)
	if err != nil {
		return err
	}

	return nil
}

// 修改页面标题的值
func modifyTitle(filePath, newVal string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var title Title
	err = json.Unmarshal(file, &title)
	if err != nil {
		return err
	}

	title.NavigationBarTitleText = newVal
	byteValue, err := json.Marshal(title)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, byteValue, 0777)
	if err != nil {
		return err
	}

	return nil
}

// 修改logo
func modifyLogo(filePath string, data []byte) error {
	_, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		return err
	}

	return nil
}

// 查询上一次配置信息 unused
func QueryLastConfig(c *Context, form XcxBaseForm) {
	storePath := conf.C.Respository.Store                                                                    // 小程序代码私库根目录
	xcxConfFileName := storePath + "/" + form.StoreId + "/" + form.VersionId + confPath + "/" + confFileName // 医院配置信息文件

	// 判断文件是否存在
	b, _ := utils.PathExists(xcxConfFileName)
	if !b {
		c.Success(nil)
		return
	}

	file, err := ioutil.ReadFile(xcxConfFileName)
	// if err != nil {
	// 	utils.Error(err)
	// 	c.Error5xx(err)
	// 	return
	// }
	// c.Success(string(file))

	var config XcxConfigForm
	err = json.Unmarshal(file, &config)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}

	c.Success(config)

}

// 公库新增版本(版本号重复会覆盖)
func AddNewVersion(c *Context, form AddXcxProjectForm) {

	versionId := form.VersionId

	destPath := conf.C.Respository.Root + "/" + versionId

	var flag bool

	// 判断目录是否存在
	b, err := utils.PathExists(destPath)
	if err != nil {
		utils.Error(err)
		c.Error5xx(err)
		return
	}
	// 不存在则创建
	if !b {
		err = os.MkdirAll(destPath, 0666)
		if err != nil {
			utils.Error(err)
			c.Error5xx(err)
			return
		}
	} else {
		// err = fmt.Errorf("当前版本已经存在,请确认版本号!")
		// c.Error5xx(err)
		// return
		flag = true
	}

	fileName := form.ProjectFile.Filename // 上传文件名

	// 创建新文件
	archive := destPath + "/" + fileName
	newFile, err := os.OpenFile(archive, os.O_RDWR|os.O_CREATE, 0777)
	defer newFile.Close()
	if err != nil {
		c.Error5xx(err)
		return
	}

	// 读取上传文件内容
	fileReader, err := form.ProjectFile.Open()
	if err != nil {
		c.Error5xx(err)
		return
	}
	fileData, err := ioutil.ReadAll(fileReader)
	if err != nil {
		c.Error5xx(err)
		return
	}

	// 写入文件
	_, err = newFile.Write(fileData)
	if err != nil {
		c.Error5xx(err)
		return
	}

	// 解压文件(压缩包)
	err = utils.Unzip(archive, destPath)
	if err != nil {
		c.Error5xx(err)
		return
	}

	// 删除.git目录
	gitPath := destPath + "/" + conf.C.Project.Name + "/.git"
	os.RemoveAll(gitPath)

	// 删除压缩包 TODO 文件占用无法删除
	// defer utils.ForceDelFile(archive)

	// go func() {
	// 	time.Sleep(60 * time.Second)
	// 	utils.ForceDelFile(archive)
	// }()

	if !flag {
		c.Success("新版本添加成功")
	} else {
		c.Success("版本添加成功(版本覆盖)")
	}

}
