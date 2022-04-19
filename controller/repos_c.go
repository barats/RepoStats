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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 转向到 代码仓库页面
//
func ReposPage(ctx *gin.Context) {

	strPage := ctx.DefaultQuery("page", strconv.Itoa(storage.DEFAULT_PAGE_NUMBER))
	strSize := ctx.DefaultQuery("size", strconv.Itoa(storage.DEFAULT_PAGE_SIZE))

	page, err := strconv.Atoi(strPage)
	if err != nil || page < storage.DEFAULT_PAGE_NUMBER {
		page = storage.DEFAULT_PAGE_NUMBER
	}

	size, err := strconv.Atoi(strSize)
	if err != nil || size < storage.DEFAULT_PAGE_SIZE || size > storage.DEFAULT_MAX_PAGE_SIZE {
		size = storage.DEFAULT_PAGE_SIZE
	}

	count, err1 := gitee_storage.FindTotalReposCount()
	repos, err2 := gitee_storage.FindPagedRepos(page, size)

	if err1 != nil || err2 != nil {
		log.Printf("error1 %s \t error2 %s", err1, err2)
		ctx.HTML(http.StatusOK, "repos.html", gin.H{
			"title":       "代码仓库列表 - RepoStats",
			"current_url": ctx.Request.URL.Path,
			"error":       "内部错误，请联系管理员",
		})
		return
	}

	ctx.HTML(http.StatusOK, "repos.html", gin.H{
		"title":        "代码仓库列表 - RepoStats",
		"current_url":  ctx.Request.URL.Path,
		"repos":        repos,
		"total_item":   count,
		"current_page": page,
		"page_size":    size,
		"first_page":   page == 1,
		"last_page":    page >= (count/size)+1,
	})
}

// 禁用、启用 代码仓库的爬取功能
//
func RepoStateChange(ctx *gin.Context) {
	strRepoID := ctx.Param("repoID")
	strType := ctx.DefaultPostForm("type", "gitee")
	strEnable := ctx.PostForm("enable")

	repoID, err := strconv.Atoi(strRepoID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "repoID 参数非法",
		})
		return
	}

	if !strings.EqualFold(strings.ToLower(strType), "gitee") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "type 参数非法",
		})
		return
	}

	enable, err := strconv.ParseBool(strEnable)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "enable 参数非法",
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
			"message": "找不到指定的仓库，repoID=" + strRepoID,
		})
		return
	}

	repo.EnableCrawl = enable
	err = gitee_storage.BulkSaveRepos([]gitee_model.Repository{repo})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "内部错误，请联系管理员！ ",
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// 删除代码仓库
//
func RepoDelete(ctx *gin.Context) {
	strRepoID := ctx.Param("repoID")
	strType := ctx.DefaultPostForm("type", "gitee")

	repoID, err := strconv.Atoi(strRepoID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "repoID 参数非法",
		})
		return
	}

	if !strings.EqualFold(strings.ToLower(strType), "gitee") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "type 参数非法",
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
			"message": "找不到指定的仓库，repoID=" + strRepoID,
		})
		return
	}

	err = gitee_storage.DeleteRepo(repoID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "内部错误，请联系管理员！ ",
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
