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
)

type Commit struct {
	RepoID  int    `db:"repo_id"`
	Sha     string `db:"sha" json:"sha"`
	HtmlUrl string `db:"html_url" json:"html_url"`

	Committer User `json:"committer"`
	Author    User `json:"author"`

	Detail struct {
		Author struct {
			User
			Date time.Time `json:"date" db:"date"`
		} `json:"author" db:"author"`
		Committer struct {
			User
			Date time.Time `json:"date" db:"date"`
		} `json:"committer" db:"committer"`
		Message string `json:"message" db:"message"`
		Tree    struct {
			Sha string `json:"sha"`
		} `json:"tree"`
	} `json:"commit" db:"commit"`
}

func (c Commit) IsNilOrEmpty() bool {
	return reflect.DeepEqual(c, Commit{})
}
