// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Dashbaord page controller
//
// Display dashbaord page html
func DashboardPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":       "仪表盘 - RepoStats",
		"current_url": ctx.Request.URL.Path,
	})
}
