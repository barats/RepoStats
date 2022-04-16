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
	"repostats/network"
	"repostats/utils"
	"testing"
)

func TestBulkSaveUsers(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	found, err := network.GetGiteeCommits("openharmony", "community")
	utils.ExitOnError(err)

	var users []gitee_model.User
	for _, f := range found {
		users = append(users, f.Author)
		users = append(users, f.Committer)
	}

	//remove duplicate users
	inResult := make(map[int]bool)
	var result []gitee_model.User
	for _, u := range users {
		if _, ok := inResult[u.ID]; !ok {
			inResult[u.ID] = true
			result = append(result, u)
		}
	}

	type args struct {
		users []gitee_model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase save user", args: args{users: result}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkSaveUsers(tt.args.users); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
