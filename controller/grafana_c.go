// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"repostats/network"
	"repostats/schedule"
	gitee_storage "repostats/storage/gitee"
	"repostats/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GrafanForm struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Account  string `json:"account"`
}

// 转向到 Grafana 配置页面
//
func GrafanaPage(ctx *gin.Context) {

	token, _ := network.RetrieveGrafanaToken()
	ctx.HTML(http.StatusOK, "grafana.html", gin.H{
		"title":       "Grafana 配置 - RepoStats",
		"current_url": ctx.Request.URL.Path,
		"grafana":     utils.GrafanaConfig,
		"token":       token,
	})
}

// 新增或修改 Grafana token
//
func GrafanaToken(ctx *gin.Context) {
	var data GrafanForm
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := network.CreateGrafanaApiToken(data.Host, data.Port, data.Account, data.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := network.RetrieveGrafanaToken()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = network.CreateDatasource(token, utils.DatabaseConifg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	datasource, err := network.RetrieveGrafanaDatasource()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = network.CreateGiteeRepostatsFolder(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	folder, err := network.RetrieveGiteeRepostatsFolder()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = network.CreateGiteeHomeDashboard(token, folder, datasource)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	rs, _ := json.Marshal(token)
	ctx.JSON(http.StatusOK, rs)
}

//创建或更新制定仓库的 Grafana 视图面板
//
func CreateOrUpdateGrafanaDashboard(ctx *gin.Context) {
	strRepoID := ctx.Param("repoID")

	repoID, err := strconv.Atoi(strRepoID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "repoID 参数非法",
		})
		return
	}

	if repoID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "repoID 参数非法",
		})
		return
	}

	repo, err := gitee_storage.FindRepoByID(repoID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "内部错误，请联系管理员！ ",
		})
		return
	}

	if repo.IsNilOrEmpty() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "找不到指定的仓库，repoID =" + strRepoID,
		})
		return
	}

	token, _ := network.RetrieveGrafanaToken()
	datasource, _ := network.RetrieveGrafanaDatasource()
	folder, _ := network.RetrieveGiteeRepostatsFolder()

	err = schedule.CreateOrUpdateGrafanaRepo(repo, token, folder, datasource)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "内部错误，请联系管理员！ " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
