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
	gitee_model "repostats/model/gitee"
	"repostats/network"
	"repostats/utils"
	"testing"
)

func TestBulkSavePullRequests(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	found, err := network.GetGiteePullRequests("openharmony", "community")
	utils.ExitOnError(err)

	for i := 0; i < len(found); i++ {
		found[i].RepoID = 10918992
	}

	if len(found) > 0 {
		var users []gitee_model.User
		for i := 0; i < len(found); i++ {
			users = append(users, found[i].User)
		}
		users = gitee_model.RemoveDuplicateUsers(users)
		err := BulkSaveUsers(users)
		utils.ExitOnError(err)
	}

	type args struct {
		prs []gitee_model.PullRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase openharmony/community", args: args{prs: found}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkSavePullRequests(tt.args.prs); (err != nil) != tt.wantErr {
				t.Errorf("BulkSavePullRequests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindPRs(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "TestCase FindPRs", want: 869, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindPRs()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPRs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("FindPRs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindPRByID(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		prID int
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{name: "TestCase TestFindPRByID", args: args{prID: 5855633}, want: 850, wantErr: false},
		{name: "TestCase TestFindPRByID", args: args{prID: 111}, want: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindPRByID(tt.args.prID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPRByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Number, tt.want) {
				t.Errorf("FindPRByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
