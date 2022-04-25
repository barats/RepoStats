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
	"repostats/network"
	"repostats/schedule"
	"repostats/storage"
	gitee_storage "repostats/storage/gitee"
	"repostats/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 添加仓库
//
// 支持添加单个仓库，也可以指定个人首页地址从而添加该用户名下所有公开仓库
func AddRepo(ctx *gin.Context) {
	repoUrl := ctx.PostForm("repo_url")
	repoType := ctx.PostForm("type")

	if utils.EmptyString(repoType) || utils.EmptyString(repoUrl) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "repo_url 或 repo_type 未指定",
		})
		return
	}

	owner, repo, err := utils.ParseGiteeRepoUrl(strings.TrimSpace(repoUrl))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !utils.EmptyString(repo) {
		//仓库链接
		found, err := network.GetGiteeRepo(owner, repo)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		found.EnableCrawl = true //默认开启抓取
		err = gitee_storage.BulkSaveRepos([]gitee_model.Repository{found})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		//个人链接 或者组织链接
		var repos []gitee_model.Repository
		repos, err := network.GetGiteeUserRepos(owner)
		if err != nil {
			repos, err = network.GetGiteeOrgRepos(owner)
		}

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if len(repos) <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "未找到任何仓库信息",
			})
			return
		}

		for i := 0; i < len(repos); i++ {
			repos[i].EnableCrawl = true //默认开启抓取
		}
		err = gitee_storage.BulkSaveRepos(repos)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, nil)
}

// 转向到 代码仓库页面
//
func ReposPage(ctx *gin.Context) {

	strPage := ctx.DefaultQuery("page", strconv.Itoa(storage.DEFAULT_PAGE_NUMBER))
	strSize := ctx.DefaultQuery("size", strconv.Itoa(storage.DEFAULT_PAGE_SIZE))
	strRepoID := ctx.DefaultQuery("id", "")
	repoName := ctx.DefaultQuery("repo_name", "")

	page, err := strconv.Atoi(strPage)
	if err != nil || page < storage.DEFAULT_PAGE_NUMBER {
		page = storage.DEFAULT_PAGE_NUMBER
	}

	size, err := strconv.Atoi(strSize)
	if err != nil || size < storage.DEFAULT_PAGE_SIZE || size > storage.DEFAULT_MAX_PAGE_SIZE {
		size = storage.DEFAULT_PAGE_SIZE
	}

	repoID, err := strconv.Atoi(strRepoID)
	if err != nil {
		repoID = 0
	}

	var count int
	var repos []gitee_model.Repository
	if repoID <= 0 {

		if utils.EmptyString(repoName) {
			count, err = gitee_storage.FindTotalReposCount()
			repos, err = gitee_storage.FindPagedRepos(page, size)

			if err != nil {
				log.Printf("error %s", err)
				ctx.HTML(http.StatusOK, "repos.html", gin.H{
					"title":       "代码仓库列表 - RepoStats",
					"current_url": ctx.Request.URL.Path,
					"error":       "内部错误，请联系管理员",
				})
				return
			}
		} else {
			//search by repo_name
			count, err = gitee_storage.FindReposCountByName(repoName)
			repos, err = gitee_storage.FindPagedReposByName(repoName, page, size)

			if err != nil {
				log.Printf("error %s", err)
				ctx.HTML(http.StatusOK, "repos.html", gin.H{
					"title":       "代码仓库列表 - RepoStats",
					"current_url": ctx.Request.URL.Path,
					"error":       "内部错误，请联系管理员",
				})
				return
			}
		}
	} else {
		r, errr := gitee_storage.FindRepoByID(repoID)
		if errr != nil {
			log.Printf("error %s", err)
			ctx.HTML(http.StatusOK, "repos.html", gin.H{
				"title":       "代码仓库列表 - RepoStats",
				"current_url": ctx.Request.URL.Path,
				"error":       "内部错误，请联系管理员",
			})
			return
		}
		repos = append(repos, r)
		count = len(repos)
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
		"repo_name":    repoName,
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

func StartToGrab(ctx *gin.Context) {
	go schedule.StarGiteeJobs(false)
	ctx.JSON(http.StatusOK, nil)
}
