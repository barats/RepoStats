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

func BulkSaveCommits(commits []gitee_model.Commit) error {
	query := `INSERT INTO gitee.commits (repo_id, sha, html_url, author_name, author_email, author_date, committer_name, 
		committer_email, committer_date, detail_message,tree)
	VALUES(:repo_id,:sha,:html_url,:commit.author.name, :commit.author.email,:commit.author.date,
		 :commit.committer.name,:commit.committer.email,:commit.committer.date,:commit.message,:commit.tree.sha)
	ON CONFLICT (repo_id,sha) DO UPDATE SET repo_id=EXCLUDED.repo_id,html_url=EXCLUDED.html_url,
	author_name=EXCLUDED.author_name,author_email=EXCLUDED.author_email,author_date=EXCLUDED.author_date,committer_name=EXCLUDED.committer_name,
	committer_email=EXCLUDED.committer_email,committer_date=EXCLUDED.committer_date,detail_message=EXCLUDED.detail_message,tree=EXCLUDED.tree`
	return storage.DbNamedExec(query, commits)
}
