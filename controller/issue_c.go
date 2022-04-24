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

	"github.com/gin-gonic/gin"
)

func IssuesPage(ctx *gin.Context) {

	strPage := ctx.DefaultQuery("page", strconv.Itoa(storage.DEFAULT_PAGE_NUMBER))
	strSize := ctx.DefaultQuery("size", strconv.Itoa(storage.DEFAULT_PAGE_SIZE))
	number := ctx.DefaultQuery("number", "")

	page, err := strconv.Atoi(strPage)
	if err != nil || page < storage.DEFAULT_PAGE_NUMBER {
		page = storage.DEFAULT_PAGE_NUMBER
	}

	size, err := strconv.Atoi(strSize)
	if err != nil || size < storage.DEFAULT_PAGE_SIZE || size > storage.DEFAULT_MAX_PAGE_SIZE {
		size = storage.DEFAULT_PAGE_SIZE
	}

	if utils.EmptyString(number) {
		iss, err1 := gitee_storage.FindPagedIssues(page, size)
		count, err2 := gitee_storage.FindTotalIssuesCount()

		if err1 != nil || err2 != nil {
			log.Printf("error %s, %s ", err1, err2)
			ctx.HTML(http.StatusOK, "issues.html", gin.H{
				"title":       "Issue 列表 - RepoStats",
				"current_url": ctx.Request.URL.Path,
				"error":       "内部错误，请联系管理员",
			})
			return
		}

		ctx.HTML(http.StatusOK, "issues.html", gin.H{
			"title":        "Issues 列表 - RepoStats",
			"current_url":  ctx.Request.URL.Path,
			"issues":       iss,
			"total_item":   count,
			"current_page": page,
			"page_size":    size,
			"first_page":   page == 1,
			"number":       number,
			"last_page":    page >= (count/size)+1,
		})

		return
	}

	iss, err1 := gitee_storage.FindPagedIssuesByNumber(number, page, size)
	count, err2 := gitee_storage.FindIssuesCountByNumber(number)

	if err1 != nil || err2 != nil {
		log.Printf("error %s , %s ", err1, err2)
		ctx.HTML(http.StatusOK, "issues.html", gin.H{
			"title":       "Issue 列表 - RepoStats",
			"current_url": ctx.Request.URL.Path,
			"error":       "内部错误，请联系管理员",
		})
		return
	}

	var issList []gitee_model.Issue
	if !iss.IsNilOrEmpty() {
		issList = append(issList, iss)
	}

	ctx.HTML(http.StatusOK, "issues.html", gin.H{
		"title":        "Issues 列表 - RepoStats",
		"current_url":  ctx.Request.URL.Path,
		"issues":       issList,
		"total_item":   count,
		"current_page": page,
		"page_size":    size,
		"first_page":   page == 1,
		"last_page":    page >= (len(issList)/size)+1,
	})
}
