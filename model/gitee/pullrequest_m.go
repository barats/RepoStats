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

type PullRequest struct {
	ID            int64     `json:"id" db:"id"`
	RepoID        int64     `json:"repo_id" db:"repo_id"`
	HTMLURL       string    `json:"html_url" db:"html_url"`
	DiffUrl       string    `json:"diff_url" db:"diff_url"`
	PatchUrl      string    `json:"patch_url" db:"patch_url"`
	Number        int64     `json:"number" db:"number"`
	State         string    `json:"state" db:"state"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     null.Time `json:"updated_at" db:"updated_at"`
	ClosedAt      null.Time `json:"closed_at" db:"closed_at"`
	MergedAt      null.Time `json:"merged_at" db:"merged_at"`
	Mergeable     bool      `json:"mergeable" db:"mergeable"`
	CanMergeCheck bool      `json:"can_merge_check" db:"can_merge_check"`
	Title         string    `json:"title" db:"title"`
	User          User      `json:"user" db:"user"`
	Head          struct {
		Label string     `json:"label" db:"label"`
		Ref   string     `json:"ref" db:"ref"`
		Sha   string     `json:"sha" db:"sha"`
		User  User       `json:"user" db:"user"`
		Repo  Repository `json:"repo" db:"repo"`
	} `json:"head" db:"head"`
}

func (pr PullRequest) IsNilOrEmpty() bool {
	return reflect.DeepEqual(pr, PullRequest{})
}
