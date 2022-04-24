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
	"net/http"
	"repostats/network"
	"repostats/utils"

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

	rs, _ := json.Marshal(token)
	ctx.JSON(http.StatusOK, rs)
}
