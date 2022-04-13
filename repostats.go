// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"repostats/controller"
	"repostats/storage"
	"repostats/utils"

	"github.com/gin-gonic/gin"
)

var (
	Version = "1.0"
	Build   = "2204111911"

	//go:embed assets/* templates/*
	FS embed.FS

	cmdConfig string
)

func main() {

	flag.StringVar(&cmdConfig, "c", "repostats.ini", "config file path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, `RepoStats version:%s build %s 
		Usage: repostats [-c config_file_path]`, Version, Build)
		flag.PrintDefaults()
	}

	flag.Parse()

	_, err := utils.InitConfig(cmdConfig)
	utils.ExitOnError(err)

	_, err = storage.InitDatabaseService()
	utils.ExitOnError(err)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	initRouter(router)

	log.Println(fmt.Sprintf("[RepoStats v%s build:%s] starts at http://localhost:%d", Version, Build, utils.RepoStatsConfig.AdminPort))
	router.Run(fmt.Sprintf("localhost:%d", utils.RepoStatsConfig.AdminPort))
}

func initRouter(router *gin.Engine) {
	sub, err := fs.Sub(FS, "assets")
	utils.ExitOnError(err)

	router.StaticFS("/assets", http.FS(sub))

	tmpl, err := template.New("").ParseFS(FS, "templates/*.html")
	utils.ExitOnError(err)

	router.GET("/login", controller.Login)
	router.POST("/login", controller.DoLogin)

	admin := router.Group("/admin", controller.AdminAuthHanlder())
	admin.POST("/logout", controller.DoLogout)
	admin.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusFound, "/admin/dashboard") })
	admin.GET("/dashboard", controller.DashboardPage)
	admin.GET("/schedule", controller.SchedulePage)
	admin.GET("/gitee", controller.GiteePage)
	admin.GET("/github", controller.GithubPage)
	admin.GET("/repos", controller.ReposPage)
	admin.GET("/grafana", controller.GrafanaPage)

	public := router.Group("/admin") //Same url path with /admin WITHOUT auth handler
	public.POST("/gitee/token", controller.GiteeTokenRetrieve)

	router.SetHTMLTemplate(tmpl)
}
