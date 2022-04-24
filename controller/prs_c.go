// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package controller

import (
	"log"
	"net/http"
	gitee_model "repostats/model/gitee"
	"repostats/storage"
	gitee_storage "repostats/storage/gitee"
	"repostats/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 转向到 Pull Request 列表页面
//
func PRsPage(ctx *gin.Context) {

	strPage := ctx.DefaultQuery("page", strconv.Itoa(storage.DEFAULT_PAGE_NUMBER))
	strSize := ctx.DefaultQuery("size", strconv.Itoa(storage.DEFAULT_PAGE_SIZE))
	strPRID := ctx.DefaultQuery("id", "")

	page, err := strconv.Atoi(strPage)
	if err != nil || page < storage.DEFAULT_PAGE_NUMBER {
		page = storage.DEFAULT_PAGE_NUMBER
	}

	size, err := strconv.Atoi(strSize)
	if err != nil || size < storage.DEFAULT_PAGE_SIZE || size > storage.DEFAULT_MAX_PAGE_SIZE {
		size = storage.DEFAULT_PAGE_SIZE
	}

	prID, err := strconv.Atoi(strPRID)
	if err != nil {
		prID = 0
	}

	var count int
	var prs []gitee_model.PullRequest
	if prID <= 0 {
		count, err = gitee_storage.FindTotalPRsCount()
		prs, err = gitee_storage.FindPagedPRs(page, size)

		if err != nil {
			log.Printf("error %s", err)
			ctx.HTML(http.StatusOK, "prs.html", gin.H{
				"title":       "Pull Request 列表 - RepoStats",
				"current_url": ctx.Request.URL.Path,
				"error":       "内部错误，请联系管理员",
			})
			return
		}
	} else {
		pr, errr := gitee_storage.FindPRByID(prID)
		if errr != nil {
			log.Printf("error %s", err)
			ctx.HTML(http.StatusOK, "prs.html", gin.H{
				"title":       "Pull Request 列表 - RepoStats",
				"current_url": ctx.Request.URL.Path,
				"error":       "内部错误，请联系管理员",
			})
			return
		}
		count = 1
		prs = append(prs, pr)
	}

	ctx.HTML(http.StatusOK, "prs.html", gin.H{
		"title":        "Pull Request 列表 - RepoStats",
		"current_url":  ctx.Request.URL.Path,
		"prs":          prs,
		"total_item":   count,
		"current_page": page,
		"page_size":    size,
		"first_page":   page == 1,
		"last_page":    page >= (count/size)+1,
	})
}

// 删除 commit
//
func PRDelete(ctx *gin.Context) {
	strID := ctx.Param("id")
	strType := ctx.DefaultPostForm("type", "gitee")

	if utils.EmptyString(strID) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "PR ID 参数非法",
		})
		return
	}

	if !strings.EqualFold(strings.ToLower(strType), "gitee") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "type 参数非法",
		})
		return
	}

	prID, err := strconv.Atoi(strID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "PR ID 参数非法",
		})
		return
	}

	found, err := gitee_storage.FindPRByID(prID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "内部错误，请联系管理员！ ",
		})
		return
	}

	if found.IsNilOrEmpty() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "找不到指定的 PullRequest, ID=" + strID,
		})
		return
	}

	err = gitee_storage.DeletePR(prID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "内部错误，请联系管理员！ ",
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
