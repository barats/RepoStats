// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package controller

import (
	"fmt"
	"net/http"
	gitee_model "repostats/model/gitee"
	"repostats/network"
	"repostats/utils"

	"github.com/gin-gonic/gin"
)

type GiteeOauthForm struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUrl  string `json:"redirect_url"`
	Code         string `json:"code"`
}

// Gitee page controller
//
// Display gitee config page html
func GiteePage(ctx *gin.Context) {

	data, err := utils.ReadRepoStatsFile(gitee_model.GITEE_TOKEN_FILE)

	ctx.HTML(http.StatusOK, "gitee.html", gin.H{
		"title":       "Gitee 配置 - RepoStats",
		"current_url": ctx.Request.URL.Path,
		"error":       err,
		"oauth_info":  string(data),
	})
}

// Get gitee token
//
// Retrieve access token from gitee oath server
func GiteeTokenRetrieve(ctx *gin.Context) {
	var data GiteeOauthForm
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	url := fmt.Sprintf(`%s?grant_type=authorization_code&code=%s&client_id=%s&redirect_uri=%s`,
		gitee_model.GITEE_OAUTH_TOKEN_URL, data.Code, data.ClientID, data.RedirectUrl)

	rc, rs, err := network.HttpPost("", url, nil, map[string]string{
		"client_secret": data.ClientSecret,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if rc != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": rs,
		})
		return
	}

	if err := utils.WriteRepoStatsFile(gitee_model.GITEE_TOKEN_FILE, []byte(rs)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, rs)
}
