// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package gitee

import (
	gitee_mode "repostats/model/gitee"
	"repostats/storage"
)

func BulkSavePullRequests(prs []gitee_mode.PullRequest) error {
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
