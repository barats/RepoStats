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

type User struct {
	ID        int       `json:"id" db:"id"`
	Login     string    `json:"login" db:"login"`
	Name      string    `json:"name" db:"name"`
	AvatarURL string    `json:"avatar_url" db:"avatar_url"`
	HTMLURL   string    `json:"html_url" db:"html_url"`
	Remark    string    `json:"remark" db:"remark"`
	Type      string    `json:"type" db:"type"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (u User) IsNilOrEmpty() bool {
	return reflect.DeepEqual(u, User{})
}

func RemoveDuplicateUsers(users []User) []User {
	keys := make(map[int]bool)
	list := []User{}

	for _, item := range users {
		if _, value := keys[item.ID]; !value {
			keys[item.ID] = true
			list = append(list, item)
		}
	}
	return list
}
