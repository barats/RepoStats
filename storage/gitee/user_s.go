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

func BulkSaveUsers(users []gitee_model.User) error {
	query := `INSERT INTO gitee.users (id, login, "name", avatar_url, html_url, remark, "type", email, created_at)
	VALUES(:id,:login,:name,:avatar_url,:html_url,:remark,:type,:email,:created_at)
	ON CONFLICT (id) DO UPDATE SET login=EXCLUDED.login,name=EXCLUDED.name,avatar_url=EXCLUDED.avatar_url,html_url=EXCLUDED.html_url,
	remark=EXCLUDED.remark,type=EXCLUDED.type,email=EXCLUDED.email,created_at=EXCLUDED.created_at`
	return storage.DbNamedExec(query, users)
}
