package main

import (
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
	"yegoo.com/yegoo-marking-publish/conf"
	"yegoo.com/yegoo-marking-publish/routes"
	"yegoo.com/yegoo-marking-publish/utils"
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
	m.Get("login", routes.Login)
	m.Post("open", binding.MultipartForm(routes.XcxConfigForm{}), routes.ValidErrorHandler, routes.OpenProject)
	m.Get("preview", binding.Bind(routes.XcxBaseForm{}), routes.Preview)
	m.Get("upload", binding.Bind(routes.XcxBaseForm{}), routes.Upload)
	//m.Get("queryLastConfig", binding.Bind(routes.XcxBaseForm{}), routes.QueryLastConfig)
	m.Post("addNewVersion", binding.MultipartForm(routes.AddXcxProjectForm{}), routes.ValidErrorHandler, routes.AddNewVersion)
	m.Run(conf.C.Server.Port)

}
