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

type Stargazer struct {
	RepoID int64     `db:"repo_id"`
	StarAt time.Time `json:"star_at" db:"star_at"`
	User   `db:"user"`
}

func (s Stargazer) isNilOrEmpty() bool {
	return reflect.DeepEqual(s, Stargazer{})
}

type Collaborator struct {
	RepoID      int64 `db:"repo_id"`
	User        `db:"user"`
	Permissions struct {
		Pull  bool `json:"pull" db:"can_pull"`
		Push  bool `json:"push" db:"can_push"`
		Admin bool `json:"admin" db:"can_admin"`
	} `db:"permissions"`
}

func (c Collaborator) IsNilOrEmpty() bool {
	return reflect.DeepEqual(c, Collaborator{})
}
