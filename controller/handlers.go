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
	"repostats/storage"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := c.Cookie("RepoStatsAdmin")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		cookie, err := c.Cookie("RepoStatsCookie")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		if len(user) <= 0 || len(cookie) <= 0 {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		found, err := storage.FindAdminByAccount(user)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		if found.IsEmpty() {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		cValue, err := AdminCookieValue(found)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		if !strings.EqualFold(cValue, cookie) {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.Next()
	} //end of func
}
