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

func BulkSaveStargazers(s []gitee_mode.Stargazers) error {
	query := `INSERT INTO gitee.stargazers (user_id, repo_id, star_at) VALUES(:user.id,:repo_id,:star_at) 
	ON CONFLICT (user_id,repo_id) DO UPDATE SET user_id=EXCLUDED.user_id,repo_id=EXCLUDED.repo_id,star_at=EXCLUDED.star_at`
	return storage.DbNamedExec(query, s)
}
