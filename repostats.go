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
	"repostats/schedule"
	"repostats/storage"
	"repostats/utils"

	"github.com/Masterminds/sprig"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed assets/* templates/*
	FS embed.FS

	cmdConfig string

	group errgroup.Group
)

func main() {

	flag.StringVar(&cmdConfig, "c", "repostats.ini", "config file path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, `RepoStats v:%s build %s 
		Usage: repostats [-c config_file_path]`, utils.Version, utils.Build)
		flag.PrintDefaults()
	}

	flag.Parse()

	_, err := utils.InitConfig(cmdConfig)
	utils.ExitOnError(err)

	utils.InitWaitingGruop()

	_, err = storage.InitDatabaseService()
	utils.ExitOnError(err)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	initRouter(router)

	group.Go(func() error {
		log.Println(fmt.Sprintf("[RepoStats v%s build:%s] starts at http://localhost:%d", utils.Version, utils.Build, utils.RepoStatsConfig.AdminPort))
		return router.Run(fmt.Sprintf(":%d", utils.RepoStatsConfig.AdminPort))
	})

	group.Go(func() error {
		log.Println(fmt.Sprintf("[RepoStats v%s build:%s] schedule started. ", utils.Version, utils.Build))
		return schedule.StartGiteeSchedule()
	})

	err = group.Wait()
	utils.ExitOnError(err)
}

func initRouter(router *gin.Engine) {
	sub, err := fs.Sub(FS, "assets")
	utils.ExitOnError(err)

	router.StaticFS("/assets", http.FS(sub))

	tmpl, err := template.New("").Funcs(sprig.FuncMap()).ParseFS(FS, "templates/*.html")
	utils.ExitOnError(err)

	router.GET("/login", controller.Login)
	router.POST("/login", controller.DoLogin)
	router.GET("/captcha/:imageId", controller.ServeCaptchaImage)
	router.POST("/captcha", controller.RequestCaptchaImage)

	admin := router.Group("/admin", controller.AdminAuthHandler())
	admin.POST("/logout", controller.DoLogout)
	admin.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusFound, "/admin/gitee") })
	admin.GET("/gitee", controller.GiteePage)

	admin.GET("/repos", controller.ReposPage)
	admin.PUT("/repos/:repoID/change_state", controller.RepoStateChange)
	admin.POST("/repos/:repoID/delete", controller.RepoDelete)
	admin.POST("/repos", controller.AddRepo)
	admin.POST("/repos/grab", controller.StartToGrab)

	admin.GET("/commits", controller.CommitsPage)
	admin.POST("/commits/:sha/delete", controller.CommitDelete)

	admin.GET("/prs", controller.PRsPage)
	admin.POST("/prs/:id/delete", controller.PRDelete)

	admin.GET("/grafana", controller.GrafanaPage)
	admin.POST("/grafana/token", controller.GrafanaToken)
	admin.PUT("/grafana/:repoID/update", controller.CreateOrUpdateGrafanaDashboard)

	admin.GET("/issues", controller.IssuesPage)

	public := router.Group("/admin") //Same url path WITHOUT admin auth handler
	public.POST("/gitee/token", controller.GiteeTokenRetrieve)

	router.SetHTMLTemplate(tmpl)
}
