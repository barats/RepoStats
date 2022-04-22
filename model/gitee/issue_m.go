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

type Issue struct {
	ID             int        `json:"id" db:"id"`
	RepoID         int64      `json:"repo_id" db:"repo_id"`
	HTMLURL        string     `json:"html_url" db:"html_url"`
	Number         string     `json:"number" db:"number"`
	State          string     `json:"state" db:"state"`
	Title          string     `json:"title" db:"title"`
	User           User       `json:"user" db:"user"`
	Repository     Repository `json:"repository" db:"repository"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      null.Time  `json:"updated_at" db:"updated_at"`
	FinishedAt     null.Time  `json:"finished_at" db:"finished_at"`
	PlanStarted_at null.Time  `json:"plan_started_at" db:"plan_started_at"`
	Comments       int        `json:"comments" db:"comments"`
	Priority       int        `json:"priority" db:"priority"`
	IssueType      string     `json:"issue_type" db:"issue_type"`
	SecurityHole   bool       `json:"security_hole" db:"security_hole"`
	IssueState     string     `json:"issue_state" db:"issue_state"`
}

func (i Issue) IsNilOrEmpty() bool {
	return reflect.DeepEqual(i, Issue{})
}
