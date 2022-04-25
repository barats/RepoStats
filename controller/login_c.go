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
	"repostats/model"
	"repostats/storage"
	"repostats/utils"
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// Login page controller
//
// Display login page html
func Login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录 - RepoStats",
	})
}

// Login action
//
// Ask for account and password to DO the login action
func DoLogin(ctx *gin.Context) {
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	captchaText := ctx.PostForm("captcha-text")
	captchaId := ctx.PostForm("captcha-id")

	if utils.EmptyString(account) || utils.EmptyString(password) || len(account) < 5 || len(password) < 8 {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - RepoStats",
			"error": "用户名或密码格式错误！",
		})
		return
	}

	if utils.EmptyString(captchaText) || utils.EmptyString(captchaId) || len(captchaText) < 6 {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - RepoStats",
			"error": "验证码格式错误!",
		})
		return
	}

	//验证码有效性验证
	if !captcha.VerifyString(captchaId, captchaText) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - RepoStats",
			"error": "验证码错误，请刷新页面再重新尝试！",
		})
		return
	}

	//用户名密码有效性验证
	loginUser, err := storage.FindAdminByAccount(account)
	if err != nil || loginUser.IsEmpty() {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - RepoStats",
			"error": "用户名或密码错误",
		})
		return
	}

	pwd, _ := storage.PasswordBase58Hash(password)
	if !strings.EqualFold(loginUser.Password, pwd) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - RepoStats",
			"error": "用户名或密码错误",
		})
		return
	}

	//Write Cookie to browser
	cValue, err := AdminCookieValue(loginUser)
	if err != nil {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - RepoStats",
			"error": "内部错误，请联系管理员",
		})
		return
	}
	ctx.SetCookie("RepoStatsAdmin", loginUser.Account, 3600, "/", "", false, true)
	ctx.SetCookie("RepoStatsCookie", cValue, 3600, "/", "", false, true)
	ctx.Redirect(http.StatusFound, "/admin/gitee")
}

// Logout action
//
// Clean cookies and redirect to login page
func DoLogout(ctx *gin.Context) {
	ctx.SetCookie("RepoStatsAdmin", "", -1, "/", "", false, true)
	ctx.SetCookie("RepoStatsCookie", "", -1, "/", "", false, true)
	ctx.Redirect(http.StatusFound, "/login")
}

func ServeCaptchaImage(c *gin.Context) {
	captcha.Server(200, 45).ServeHTTP(c.Writer, c.Request)
}

func RequestCaptchaImage(c *gin.Context) {
	imageId := captcha.New()
	c.JSON(http.StatusOK, gin.H{
		"result": imageId,
	})
}

func AdminCookieValue(user model.Admin) (string, error) {
	var result string
	data, err := utils.Sha256Of(user.Account + "a=" + user.Password + "=e" + strconv.Itoa(user.ID))
	if err != nil {
		log.Println(err)
		return result, err
	}
	return utils.Base58Encode(data), nil
}
