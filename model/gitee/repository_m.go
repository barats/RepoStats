// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package gitee

import (
	"reflect"
	"time"

	"gopkg.in/guregu/null.v4"
)

type Repository struct {
	ID              int       `json:"id" db:"id"`
	FullName        string    `json:"full_name" db:"full_name"`
	HumanName       string    `json:"human_name" db:"human_name"`
	Path            string    `json:"path" db:"path"`
	Name            string    `json:"name" db:"name"`
	URL             string    `json:"url" db:"url"`
	Owner           User      `json:"owner" db:"owner"`
	Assigner        User      `json:"assigner" db:"assigner"`
	Description     string    `json:"description" db:"description"`
	HTMLURL         string    `json:"html_url" db:"html_url"`
	SSHURL          string    `json:"ssh_url" db:"ssh_url"`
	Fork            bool      `json:"fork" db:"forked_repo"`
	DefaultBranch   string    `json:"default_branch" db:"default_branch"`
	ForksCount      int       `json:"forks_count" db:"forks_count"`
	StargazersCount int       `json:"stargazers_count" db:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count" db:"watchers_count"`
	License         string    `json:"license" db:"license"`
	PushedAt        null.Time `json:"pushed_at" db:"pushed_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       null.Time `json:"updated_at" db:"updated_at"`
	EnableCrawl     bool      `db:"enable_crawl"`
}

func (r Repository) IsNilOrEmpty() bool {
	return reflect.DeepEqual(r, Repository{})
}
