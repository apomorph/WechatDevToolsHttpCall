package main

import (
	"github.com/bluecatlee/WechatDevToolsHttpCall/conf"
	"github.com/bluecatlee/WechatDevToolsHttpCall/routes"
	"github.com/bluecatlee/WechatDevToolsHttpCall/utils"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

func main() {
	// Init
	if err := conf.Init(); err != nil {
		utils.Error(err)
		return
	}

	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(cache.Cacher())
	m.Use(session.Sessioner())
	m.Use(routes.Contexter())
	m.Get("midware/login", routes.AccessLimiter, routes.Login)
	m.Get("midware/open", routes.AccessLimiter, binding.Bind(routes.XcxConfigForm{}), routes.OpenProject)
	m.Get("midware/preview", routes.AccessLimiter, binding.Bind(routes.XcxBaseForm{}), routes.Preview)
	m.Get("midware/upload", routes.AccessLimiter, binding.Bind(routes.XcxBaseForm{}), routes.Upload)
	//m.Get("midware/queryLastConfig", binding.Bind(routes.XcxBaseForm{}), routes.QueryLastConfig)
	m.Post("manager/addNewVersion", binding.MultipartForm(routes.AddXcxProjectForm{}), routes.ValidErrorHandler, routes.AddNewVersion)
	m.Get("manager/invalidsession", routes.InvalidSession)
	m.Run(conf.C.Server.Port)

}
