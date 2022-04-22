// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package gitee

import (
	gitee_model "repostats/model/gitee"
	"repostats/storage"
)

var prQueryPrefix = `SELECT pr.user_id AS "user.id", pr.head_label AS "head.label", pr.head_ref AS "head.ref", pr.head_sha AS "head.sha", 
pr.head_user_id AS "head.user.id", pr.head_repo_id AS "head.repo.id",pr.* FROM gitee.pull_requests pr `

func BulkSavePullRequests(prs []gitee_model.PullRequest) error {
	query := `INSERT INTO gitee.pull_requests (id, repo_id, user_id, html_url, diff_url, patch_url, "number", 
	state, created_at, updated_at, closed_at, merged_at, mergeable, can_merge_check, title, head_label, head_ref, 
	head_sha, head_user_id, head_repo_id)
	VALUES(:id, :repo_id,:user.id,:html_url,:diff_url,:patch_url,:number,:state,:created_at,:updated_at,:closed_at,
	:merged_at, :mergeable, :can_merge_check, :title, :head.label, :head.ref, :head.sha,:head.user.id,:head.repo.id)
	ON CONFLICT (id,repo_id) DO UPDATE SET user_id=EXCLUDED.user_id, html_url=EXCLUDED.html_url,diff_url=EXCLUDED.diff_url,
	patch_url=EXCLUDED.patch_url,number=EXCLUDED.number,state=EXCLUDED.state,created_at=EXCLUDED.created_at,updated_at=EXCLUDED.updated_at,
	closed_at=EXCLUDED.closed_at,merged_at=EXCLUDED.merged_at,mergeable=EXCLUDED.mergeable,can_merge_check=EXCLUDED.can_merge_check,title=EXCLUDED.title,
	head_label=EXCLUDED.head_label,head_ref=EXCLUDED.head_ref,head_sha=EXCLUDED.head_sha,head_user_id=EXCLUDED.head_user_id,head_repo_id=EXCLUDED.head_repo_id`
	return storage.DbNamedExec(query, prs)
}

func FindTotalPRsCount() (int, error) {
	var count int
	query := `SELECT count(pr.id) FROM gitee.pull_requests pr`
	return count, storage.DbGet(query, &count)
}

func FindPagedPRs(page, size int) ([]gitee_model.PullRequest, error) {
	if page < 1 {
		page = 1
	}
	prs := []gitee_model.PullRequest{}
	query := prQueryPrefix + ` ORDER BY pr.created_at DESC LIMIT $1 OFFSET $2`
	offset := (page - 1) * size
	return prs, storage.DbSelect(query, &prs, size, offset)
}

func FindPRByID(prID int) (gitee_model.PullRequest, error) {
	found := gitee_model.PullRequest{}
	err := storage.DbGet(prQueryPrefix+` WHERE pr.id = $1`, &found, prID)
	return found, err
}

func FindPRs() ([]gitee_model.PullRequest, error) {
	found := []gitee_model.PullRequest{}
	err := storage.DbSelect(prQueryPrefix+` ORDER BY pr.created_at DESC`, &found)
	return found, err
}

func FindPRsByRepoID(repoID int) ([]gitee_model.PullRequest, error) {
	found := []gitee_model.PullRequest{}
	err := storage.DbSelect(prQueryPrefix+` WHERE pr.repo_id = $1 ORDER BY pr.created_at DESC`, &found, repoID)
	return found, err
}

func DeletePR(prID int) error {
	query := `DELETE FROM gitee.pull_requests WHERE id = $1`
	return storage.DbExec(query, prID)
}
