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
	"repostats/network"
	"repostats/utils"
	"testing"
)

func TestBulkSaveStargazers(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	found1, err := network.GetGiteeStargazers("barat", "ohurlshortener")
	utils.ExitOnError(err)

	//TEST ONLY
	for i := 0; i < len(found1); i++ {
		found1[i].RepoID = 21133399
	}

	found2, err := network.GetGiteeStargazers("openharmony", "communication_dsoftbus")
	utils.ExitOnError(err)

	//TEST ONLY
	for i := 0; i < len(found2); i++ {
		found2[i].RepoID = 16184960
	}

	type args struct {
		s []gitee_mode.Stargazer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase barat", args: args{s: found1}, wantErr: false},
		{name: "TestCase openharmony", args: args{s: found2}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkSaveStargazers(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveStargazers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBulkSaveCollaborators(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	found1, err := network.GetGiteeCollaborators("barat", "ohurlshortener")
	utils.ExitOnError(err)

	//TEST ONLY
	for i := 0; i < len(found1); i++ {
		found1[i].RepoID = 21133399
	}

	found2, err := network.GetGiteeCollaborators("openharmony", "communication_dsoftbus")
	utils.ExitOnError(err)

	//TEST ONLY
	for i := 0; i < len(found2); i++ {
		found2[i].RepoID = 16184960
	}

	type args struct {
		c []gitee_mode.Collaborator
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Testcase barat", args: args{c: found1}, wantErr: false},
		{name: "Testcase openharmony", args: args{c: found2}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkSaveCollaborators(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveCollaborators() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
