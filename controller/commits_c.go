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

// 转向到 Commit 列表页面
//
func CommitsPage(ctx *gin.Context) {

	strPage := ctx.DefaultQuery("page", strconv.Itoa(storage.DEFAULT_PAGE_NUMBER))
	strSize := ctx.DefaultQuery("size", strconv.Itoa(storage.DEFAULT_PAGE_SIZE))
	sha := ctx.DefaultQuery("sha", "")
	authorEmail := ctx.DefaultQuery("author_email", "")
	committerEmail := ctx.DefaultQuery("committer_email", "")

	page, err := strconv.Atoi(strPage)
	if err != nil || page < storage.DEFAULT_PAGE_NUMBER {
		page = storage.DEFAULT_PAGE_NUMBER
		err = nil
	}

	size, err := strconv.Atoi(strSize)
	if err != nil || size < storage.DEFAULT_PAGE_SIZE || size > storage.DEFAULT_MAX_PAGE_SIZE {
		size = storage.DEFAULT_PAGE_SIZE
		err = nil
	}

	if utils.EmptyString(sha) && utils.EmptyString(authorEmail) && utils.EmptyString(committerEmail) {

		count, err := gitee_storage.FindTotalCommitsCount()
		commits, err := gitee_storage.FindPagedCommits(page, size)

		if err != nil {
			log.Printf("error %s ", err)
			ctx.HTML(http.StatusOK, "commits.html", gin.H{
				"title":       "Commit 列表 - RepoStats",
				"current_url": ctx.Request.URL.Path,
				"error":       "内部错误，请联系管理员",
			})
			return
		}

		ctx.HTML(http.StatusOK, "commits.html", gin.H{
			"title":        "Commit 列表 - RepoStats",
			"current_url":  ctx.Request.URL.Path,
			"commits":      commits,
			"total_item":   count,
			"current_page": page,
			"page_size":    size,
			"first_page":   page == 1,
			"last_page":    page >= (count/size)+1,
		})

		return
	} //end of if

	if !utils.EmptyString(sha) {
		found, err := gitee_storage.FindCommitBySha(sha)
		if err != nil {
			log.Printf("err %s", err)
			ctx.HTML(http.StatusOK, "commits.html", gin.H{
				"title":       "Commit 列表 - RepoStats",
				"current_url": ctx.Request.URL.Path,
				"error":       "内部错误，请联系管理员",
			})
			return
		}

		var commits []gitee_model.Commit
		if !found.IsNilOrEmpty() {
			commits = append(commits, found)
		}

		ctx.HTML(http.StatusOK, "commits.html", gin.H{
			"title":        "Commit 列表 - RepoStats",
			"current_url":  ctx.Request.URL.Path,
			"commits":      commits,
			"total_item":   len(commits),
			"current_page": page,
			"page_size":    size,
			"first_page":   page == 1,
			"last_page":    page >= (len(commits)/size)+1,
			"sha":          sha,
		})

		return
	}

	if !utils.EmptyString(authorEmail) {
		count, err := gitee_storage.FindCommitsCountByAuthorEmail(authorEmail)
		commits, err := gitee_storage.FindPagedCommitsByAuthorEmail(authorEmail, page, size)

		if err != nil {
			log.Printf("err %s", err)
			ctx.HTML(http.StatusOK, "commits.html", gin.H{
				"title":       "Commit 列表 - RepoStats",
				"current_url": ctx.Request.URL.Path,
				"error":       "内部错误，请联系管理员",
			})
			return
		}

		ctx.HTML(http.StatusOK, "commits.html", gin.H{
			"title":        "Commit 列表 - RepoStats",
			"current_url":  ctx.Request.URL.Path,
			"commits":      commits,
			"total_item":   count,
			"current_page": page,
			"page_size":    size,
			"first_page":   page == 1,
			"last_page":    page >= (count/size)+1,
			"author_email": authorEmail,
		})

		return
	}

	if !utils.EmptyString(committerEmail) {

		count, err := gitee_storage.FindCommitsCountByCommitterEmail(committerEmail)
		commits, err := gitee_storage.FindPagedCommitsByCommitterEmail(committerEmail, page, size)

		if err != nil {
			log.Printf("err %s", err)
			ctx.HTML(http.StatusOK, "commits.html", gin.H{
				"title":       "Commit 列表 - RepoStats",
				"current_url": ctx.Request.URL.Path,
				"error":       "内部错误，请联系管理员",
			})
			return
		}

		ctx.HTML(http.StatusOK, "commits.html", gin.H{
			"title":           "Commit 列表 - RepoStats",
			"current_url":     ctx.Request.URL.Path,
			"commits":         commits,
			"total_item":      count,
			"current_page":    page,
			"page_size":       size,
			"first_page":      page == 1,
			"last_page":       page >= (count/size)+1,
			"committer_email": committerEmail,
		})

		return
	}
}

// 删除 commit
//
func CommitDelete(ctx *gin.Context) {
	sha := ctx.Param("sha")
	strType := ctx.DefaultPostForm("type", "gitee")

	if utils.EmptyString(sha) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "sha 参数非法",
		})
		return
	}

	if !strings.EqualFold(strings.ToLower(strType), "gitee") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "type 参数非法",
		})
		return
	}

	found, err := gitee_storage.FindCommitBySha(sha)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "内部错误，请联系管理员！ ",
		})
		return
	}

	if found.IsNilOrEmpty() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "找不到指定的 Commit，sha=" + sha,
		})
		return
	}

	err = gitee_storage.DeleteCommitBySha(sha)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "内部错误，请联系管理员！ ",
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
