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

func BulkSaveStargazers(s []gitee_mode.Stargazer) error {
	query := `INSERT INTO gitee.stargazers (user_id, repo_id, star_at) VALUES(:user.id,:repo_id,:star_at) 
	ON CONFLICT (user_id,repo_id) DO UPDATE SET user_id=EXCLUDED.user_id,repo_id=EXCLUDED.repo_id,star_at=EXCLUDED.star_at`
	return storage.DbNamedExec(query, s)
}

func BulkSaveCollaborators(c []gitee_mode.Collaborator) error {
	query := `INSERT INTO gitee.collaborators (user_id, repo_id, can_pull, can_push, can_admin)
	VALUES(:user.id,:repo_id,:permissions.can_pull,:permissions.can_push,:permissions.can_admin)
	ON CONFLICT (user_id,repo_id) DO UPDATE SET user_id=EXCLUDED.user_id,repo_id=EXCLUDED.repo_id,
	can_pull=EXCLUDED.can_pull,can_push=EXCLUDED.can_push,can_admin=EXCLUDED.can_admin`
	return storage.DbNamedExec(query, c)
}
