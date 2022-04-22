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
	"repostats/storage"
	"repostats/utils"
	"testing"
)

func testSetup(t *testing.T) {
	_, err := utils.InitConfig("../../repostats.ini")
	if err != nil {
		t.Error(err)
		utils.ExitOnError(err)
	}

	_, err = storage.InitDatabaseService()
	if err != nil {
		t.Error(err)
		utils.ExitOnError(err)
	}
}

func testTeardown(t *testing.T) {
	storage.DbClose()
}

func TestBulkSaveCommits(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	var users []gitee_model.User

	found1, err := network.GetGiteeCommits("openharmony", "community")
	utils.ExitOnError(err)

	for i := 0; i < len(found1); i++ {
		found1[i].RepoID = 10918992
		users = append(users, found1[i].Author, found1[i].Committer)
	}

	found2, err := network.GetGiteeCommits("barat", "ohurlshortener")
	utils.ExitOnError(err)
	for i := 0; i < len(found2); i++ {
		found2[i].RepoID = 21133399
		users = append(users, found2[i].Author, found2[i].Committer)
	}

	if len(users) > 0 {
		users = gitee_model.RemoveDuplicateUsers(users)
		err := BulkSaveUsers(users)
		utils.ExitOnError(err)
	}

	type args struct {
		commits []gitee_model.Commit
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase1", args: args{commits: found1}, wantErr: false},
		{name: "TestCase2", args: args{commits: found2}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkSaveCommits(tt.args.commits); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveCommits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindCommits(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "TestCase1", want: 597, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindCommits()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindCommits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("FindCommits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindCommitsByRepoID(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		repoID int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestCase 1", args: args{repoID: 111}, want: 0, wantErr: false},
		{name: "TestCase 2", args: args{repoID: 10918992}, want: 597, wantErr: false}, // 597
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindCommitsByRepoID(tt.args.repoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindCommitsByRepoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("FindCommitsByRepoID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindCommitBySha(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	type args struct {
		sha string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "TestCase 1", args: args{sha: "adsfas"}, want: "", wantErr: false},
		{name: "TestCase 2", args: args{sha: "dfb6fb1514ceb065ecd18c71a4c38ecb4b497b1c"}, want: "dfb6fb1514ceb065ecd18c71a4c38ecb4b497b1c", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindCommitBySha(tt.args.sha)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindCommitBySha() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Sha, tt.want) {
				t.Errorf("FindCommitBySha() = %v, want %v", got, tt.want)
			}
		})
	}
}
