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

func BulkSaveIssues(iss []gitee_model.Issue) error {
	query := `INSERT INTO gitee.issues (id, html_url, "number", state, title, user_id, repo_id, finished_at, created_at, 
				updated_at, plan_started_at, "comments", priority, issue_type, issue_state, security_hole)
				VALUES(:id,:html_url,:number,:state,:title,:user.id,:repository.id,:finished_at,:created_at,
				:updated_at,:plan_started_at,:comments,:priority,:issue_type, :issue_state, :security_hole)
				ON CONFLICT (id) DO UPDATE SET id=EXCLUDED.id,html_url=EXCLUDED.html_url,number=EXCLUDED.number,
				state=EXCLUDED.state,title=EXCLUDED.title,user_id=EXCLUDED.user_id,repo_id=EXCLUDED.repo_id,
				finished_at=EXCLUDED.finished_at,created_at=EXCLUDED.created_at, updated_at=EXCLUDED.updated_at,
				plan_started_at=EXCLUDED.plan_started_at,comments=EXCLUDED.comments,priority=EXCLUDED.priority,
				issue_type=EXCLUDED.issue_type,issue_state=EXCLUDED.issue_state,security_hole=EXCLUDED.security_hole`
	return storage.DbNamedExec(query, iss)
}
